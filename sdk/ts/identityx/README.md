# IdentityX SDK — TypeScript

Framework-agnostic client for the TrieOH IdentityX API: the generated API
client, the authentication services, the token layout and the shared types. It
ships no components or hooks — those live in the separate
`@trieoh/identityx-sdk-ts-react` and `@trieoh/identityx-sdk-ts-solid` packages.

**What this package owns:** the IdentityX contract (endpoints, `{code, data,
message, error_id}` envelopes, claims decoding), the token layout, and the
refresh policy (proactive when the access token is about to expire, reactive on
a 401, with retry).

**What it deliberately does not own:** *where* the session lives and *how*
calls travel. Both are injectable — that is what lets the same SDK serve a plain
browser app and an app whose tokens never reach the browser.

## Installation

```bash
pnpm add @trieoh/identityx-sdk-ts
```

## Configuration

The SDK resolves the project id from `VITE_TRIEOH_AUTH_PROJECT_ID`,
`NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID` or `PUBLIC_TRIEOH_AUTH_PROJECT_ID`, and the
API key from `TRIEOH_AUTH_API_KEY` (server-side only). `BASE_URL` defaults to
`https://api.trieauth.trieoh.com`.

To set them explicitly — this wins over the environment:

```ts
import { configure } from "@trieoh/identityx-sdk-ts";

configure({
  PROJECT_ID: "your-project-id",
  BASE_URL: "https://identityx.example.com",
});
```

## Where the session lives

|              | Tokens stored in                                                                                  | Who talks to IdentityX              |
| ------------ | ------------------------------------------------------------------------------------------------- | ----------------------------------- |
| **Direct**   | `sessionStorage` (access token, per tab) and `localStorage` (refresh token + expiries), via `browserTokenStorage` | the browser, with a `Bearer` header |
| **BFF**      | an HttpOnly cookie on your server; the browser holds nothing but an opaque signed value            | your server                         |

Both modes expose the exact same `AuthService`; switching is a wiring change,
not a rewrite.

## Direct usage

### The auth service

```ts
import { Api, createAuthService, configure } from "@trieoh/identityx-sdk-ts";

configure({ PROJECT_ID: "your-project-id" });

const auth = createAuthService(new Api());

const login = await auth.login("user@example.com", "password");
if (login.success) console.log(login.data.access_token);

const profile = await auth.getActorProfile("actor-id");
await auth.logout();
```

### Your own API calls

`createFetcher` returns the same client the services use, so every request goes
through the interceptor and gets the `Bearer` header plus the refresh handling
for free:

```ts
import { createFetcher } from "@trieoh/identityx-sdk-ts";

const api = createFetcher({ baseURL: "https://api.example.com" });

const events = await api.get<Event[]>("/events/joined");
```

### Framework bindings

`@trieoh/identityx-sdk-ts-react` and `@trieoh/identityx-sdk-ts-solid` wrap this
package in an `<AuthProvider>` + `useAuth()`. **Without** an `adapter` prop they
run in direct mode and the tokens stay in the browser:

```tsx
<AuthProvider baseURL={import.meta.env.VITE_AUTH_API_URL} projectId={projectId}>
  <App />
</AuthProvider>
```

### On a server: inject the session and the transport

There is no `sessionStorage` on a server, so the browser storages quietly no-op
and every write is lost. Give each request its own store and, if you want
tracing or retries, your own `fetch`:

```ts
import { Api, createMemoryStorage, createTokenStore } from "@trieoh/identityx-sdk-ts";

// one per request — a module-level store would leak sessions between users
const tokenStore = createTokenStore({
  access: createMemoryStorage(),
  persistent: createMemoryStorage(),
});

const api = new Api(baseURL, undefined, undefined, {
  tokenStore,
  fetch: (input, init) => fetch(input, withTrace(init)),
});

const auth = createAuthService(api);
```

## BFF usage (tokens never reach the browser)

`@trieoh/front-core/auth/bff` owns the server side: the cookie session, the
refresh, the proxy guard, and a single `POST /auth/bff {op, input}` dispatcher.
This SDK keeps owning the contract; `front-core` only decides where the tokens
live and how calls travel. The browser side is a *neutral* adapter — the same
`AuthService` shape, but every call is one `transport.call(op, input)`.

### Server — TanStack Start

```ts
import { createServerFn } from "@tanstack/react-start";
import {
  createTanStackIdentityXBff,
  type ServerProxyRequest,
} from "@trieoh/front-core-react/server";

const bff = createTanStackIdentityXBff({
  identityX: { baseURL: env.AUTH_API_URL, projectId: env.AUTH_PROJECT_ID },
  apiBaseURL: env.API_URL,
  session: { password: env.AUTH_SESSION_PASSWORD, name: "my-app-auth" },
});

export const loginServerFn = createServerFn({ method: "POST" })
  .handler(({ data }: { data: { email: string; password: string } }) =>
    bff.login(data.email, data.password));

export const logoutServerFn = createServerFn({ method: "POST" })
  .handler(() => bff.logout());

export const restoreServerFn = createServerFn({ method: "GET" })
  .handler(() => bff.restore());

// Everything else (the app's own API included) goes through here.
export const requestServerFn = createServerFn({ method: "POST" })
  .handler(({ data }: { data: ServerProxyRequest }) => bff.request(data));
```

`session.password` must be a stable secret of at least 32 characters.

### Client

```ts
import { configureApiClient, createOrvalTransport } from "@trieoh/api-client";
import { createTanStackIdentityXIntegration } from "@trieoh/front-core-react";

const { authAdapter, authFetcher } = createTanStackIdentityXIntegration(
  {
    login: loginServerFn,
    logout: logoutServerFn,
    refresh: refreshServerFn,
    restore: restoreServerFn,
    request: requestServerFn,
  },
  { projectId },
);

configureApiClient({ baseURL: "", transport: createOrvalTransport(authFetcher) });

<AuthProvider adapter={authAdapter} projectId={projectId}>
  <App />
</AuthProvider>;
```

Passing `adapter` is what switches the provider to BFF mode: `restoreSession()`
asks the server, and `isAuthenticated` comes from the cookie instead of from a
token in the browser.

### Server — any other runtime

`createBffHandler` returns a plain `(request) => Response`, and
`createBffAuthAdapter` / `createBffProxyFetchers` accept any `BffTransport`:

```ts
// server
import { createBffHandler } from "@trieoh/front-core/auth/bff/handler";

const handleBffRequest = createBffHandler({
  identityX: { baseURL: AUTH_API_URL, projectId: AUTH_PROJECT_ID },
  apiBaseURL: API_URL,
  session: { password: AUTH_SESSION_PASSWORD, name: "my-app-auth" },
});

// client
const transport = {
  call: async <T>(op: string, input?: unknown): Promise<T> => {
    const response = await fetch("/auth/bff", {
      method: "POST",
      credentials: "same-origin",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ op, input }),
    });
    return response.json() as Promise<T>;
  },
};
```

Operations: `isSetupDone`, `setup`, `login`, `loginWithProvider`,
`completeProviderLogin`, `register`, `logout`, `refresh`, `restore`,
`introspect`, `request`.

## Escape hatches

Everything the SDK would otherwise hard-wire is a parameter:

| Need                                   | Use                                                                          |
| -------------------------------------- | ---------------------------------------------------------------------------- |
| Keep the session outside the browser   | `ApiClientConfig.tokenStore` — any `TokenStore`                                |
| Store it somewhere custom              | `createTokenStore({ access, persistent })` with your own `StorageAdapter`      |
| In-memory session (tests, SSR)         | `createMemoryStorage({ trieoh_refresh_token: "..." })`                        |
| Wrap every IdentityX call              | `ApiClientConfig.fetch` — the refresh call included                            |
| Skip auth for one request              | `api.get(path, { requiresAuth: false })`                                       |
| Skip the proactive refresh for one call | `api.post(path, body, { skipRefresh: true })`                                 |

In BFF mode, `front-core` already emits browser spans (`withSpan`) and forwards
the `traceparent` header to the server, which forwards it to IdentityX.

## Service reference

| Method                                                          | Auth             |
| --------------------------------------------------------------- | ---------------- |
| `isSetupDone()`                                                 | no               |
| `setup(email, password)`                                        | no               |
| `login(email, password)`                                        | no               |
| `loginWithProvider(provider)` / `completeProviderLogin(...)`      | no               |
| `register(email, password)`                                     | no               |
| `logout({ forceLogout? })`                                      | yes              |
| `refresh()`                                                     | refresh token    |
| `profile()`                                                     | local only       |
| `sendForgotPassword`, `resetPassword`, `verifyEmail`, `resendVerifyEmail` | no      |
| `introspect(apiKey?)`                                           | yes, or `X-API-KEY` |
| `getOAuthProviders(projectId?)`                                 | no               |
| `getProjectProfile` / `getPlatformProfile` / `getActorProfile`   | yes              |
| `upsertActorProfile` / `upsertPlatformProfile` / `upsertProjectProfile` | yes       |
| `getProfileByHandle(handle)`                                    | no               |
| `getProfileSchema` / `upsertProfileSchema`                       | yes              |
| `health()`                                                      | no               |

## Behaviour worth knowing

- **Refresh is single-flight** per `Api` instance, and in the BFF it is also
  shared between concurrent requests carrying the same refresh token. The
  refresh token is single-use: a second concurrent call would log the user out.
- **Refresh happens 30s before expiry** (`isTokenExpiringSoon(30)`), and again
  reactively if the API answers 401 — one retry, then the error surfaces.
- **A 4xx from `/auth/refresh` clears the session**; a 5xx or a network failure
  keeps it, so a flaky IdentityX does not log everyone out. The BFF never clears
  a session it cannot vouch for.
- Failures come back as `{ success: false, code, error_id, message }`; check
  `success` instead of catching.
