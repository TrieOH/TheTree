# IdentityX SDK - Typescript

SDK for integrating with the TrieOH authentication ecosystem.

## Installation

```bash
npm install @trieoh/identityx-sdk-ts
# or
yarn add @trieoh/identityx-sdk-ts
# or
bun add @trieoh/identityx-sdk-ts
```

## Configuration (Vite / React)

To use the SDK in a React project (Vite, Next.js, or CRA), wrap your application with `AuthProvider`.

### Option 1: Environment Variables (Recommended)

The SDK automatically looks for these variables:

- `VITE_TRIEOH_AUTH_PROJECT_ID` (Vite)
- `NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID` (Next.js)
- `PUBLIC_TRIEOH_AUTH_PROJECT_ID` (General)

```tsx
import { AuthProvider } from '@trieoh/identityx-sdk-ts/react';

function App() {
  return (
    <AuthProvider>
      <YourRoutes />
    </AuthProvider>
  );
}
```

### Option 2: Passing via Props

Useful if you load the project ID dynamically or want to avoid environment issues.

```tsx
<AuthProvider projectId="your-project-id-here">
  <YourApp />
</AuthProvider>
```

### Option 3: Global Configuration via Code

```tsx
import { configure } from '@trieoh/identityx-sdk-ts';

configure({
  PROJECT_ID: 'your-id',
  BASE_URL: 'https://your-api.com'
});
```

## Components

The SDK provides ready-to-use components:

```tsx
import { SignIn, SignUp } from '@trieoh/identityx-sdk-ts/react';

// Example usage
const LoginPage = () => <SignIn />;
const RegisterPage = () => <SignUp />;
```

## Hooks

You can access the authentication state anywhere in your application:

```tsx
import { useAuth } from '@trieoh/identityx-sdk-ts/react';

function Header() {
  const { isAuthenticated, auth } = useAuth();

  return (
    <nav>
      {isAuthenticated ? (
        <button onClick={() => auth.logout()}>Logout</button>
      ) : (
        <span>Not logged in</span>
      )}
    </nav>
  );
}
```

## Server-managed sessions

The default provider stores tokens in browser storage for compatibility with
client-only applications. Applications with a server runtime can replace the
auth transport and keep tokens out of browser JavaScript:

```tsx
import {
  AuthProvider,
  type AuthProviderAdapter,
} from "@trieoh/identityx-sdk-ts/react";
import { createServerAuthService, restoreServerSession } from "./auth.server";

const adapter: AuthProviderAdapter = {
  restoreSession: restoreServerSession,
  createAuth: ({ callbacks, setAuthenticated }) =>
    createServerAuthService({ callbacks, setAuthenticated }),
};

export function App() {
  return (
    <AuthProvider adapter={adapter}>
      <YourRoutes />
    </AuthProvider>
  );
}
```

For TanStack Start, `restoreServerSession` and the methods returned by
`createServerAuthService` can call server functions. Login and refresh should
happen on the server, and their IdentityX token responses must not be returned
to the browser. Prefer an opaque session cookie backed by server-side storage:

- `HttpOnly`
- `Secure` in production
- `SameSite=Lax` (or stricter where possible)
- `Path=/`

State-changing server functions must also validate the request origin or use a
CSRF token. The adapter's `setAuthenticated` callback updates the React auth
state after server-side login/logout; `callbacks` preserves the standard
`AuthProvider` lifecycle callbacks.

---
