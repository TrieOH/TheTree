# IdentityX BFF for TanStack Start

This module implements server-managed IdentityX sessions specifically for
TanStack Start. Framework-specific integrations for other runtimes should live
in sibling directories under `auth/`.

## Server configuration

The consumer provides infrastructure configuration, not IdentityX endpoint
details. Login, logout and refresh paths, methods and headers are internal.

```ts
import { createTanStackIdentityXBff } from "@trieoh/front-core/auth/tanstack/server";

const bff = createTanStackIdentityXBff({
  identityX: {
    baseURL: env.VITE_API_URL,
    projectId: env.VITE_TRIEOH_AUTH_PROJECT_ID,
  },
  apiBaseURL: env.VITE_API_URL,
  session: {
    password: env.AUTH_SESSION_PASSWORD,
    name: "my-app-auth",
    secure: import.meta.env.PROD,
  },
});
```

Expose the domain operations through application-owned server functions:

```ts
export const loginServerFn = createServerFn({ method: "POST" })
  .validator(loginSchema)
  .handler(({ data }) => bff.login(data.email, data.password));

export const restoreServerFn = createServerFn({ method: "GET" })
  .handler(() => bff.restore());

export const setupServerFn = createServerFn({ method: "POST" })
  .validator(loginSchema)
  .handler(({ data }) => bff.setup(data.email, data.password));

export const proxyServerFn = createServerFn({ method: "POST" })
  .validator(proxySchema)
  .handler(({ data }) => bff.request(data));
```

The browser never supplies a base URL. The BFF always targets `apiBaseURL` and
also rejects absolute/network-path URLs and unsupported HTTP methods. The
IdentityX and API URLs may be identical; they are separate configuration fields
only because IdentityX operations have fixed, library-owned contracts.

## Client configuration

```ts
import {
  createTanStackIdentityXAuthProviderAdapter,
  createTanStackServerProxyFetchers,
} from "@trieoh/front-core/auth/tanstack/client";

export const authAdapter = createTanStackIdentityXAuthProviderAdapter({
  setup: setupServerFn,
  login: loginServerFn,
  logout: () => logoutServerFn(),
  refresh: () => refreshServerFn(),
  restore: () => restoreServerFn(),
});

export const { authFetcher, authQueryFetcher } =
  createTanStackServerProxyFetchers(proxyServerFn);
```

Only `login`, `logout`, `refresh`, and `restore` are required by the adapter.
Pass `isSetupDone`, `setup`, `register`, `loginWithProvider`, and
`completeProviderLogin` when that frontend uses those flows. Their IdentityX
URLs, methods, and token handling remain inside this package.

`AUTH_SESSION_PASSWORD` must be a stable random secret of at least 32
characters. Store it in the deployment platform's secret manager. Rotating it
invalidates all active application sessions.
