# IdentityX SDK — React

React bindings and ready-to-use authentication components for IdentityX.
The framework-agnostic client is provided by the peer package
`@trieoh/identityx-sdk-ts`.

## Installation

```bash
npm install @trieoh/identityx-sdk-ts-react @trieoh/identityx-sdk-ts
```

## Provider

```tsx
import { AuthProvider } from "@trieoh/identityx-sdk-ts-react";

export function App() {
  return (
    <AuthProvider>
      <YourRoutes />
    </AuthProvider>
  );
}
```

The provider reads `VITE_TRIEOH_AUTH_PROJECT_ID`,
`NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID` or `PUBLIC_TRIEOH_AUTH_PROJECT_ID`.
Configuration can also be set through the base SDK:

```ts
import { configure } from "@trieoh/identityx-sdk-ts";

configure({ PROJECT_ID: "your-project-id" });
```

## Components and hooks

```tsx
import {
  ModernAuth,
  SignIn,
  SignUp,
  useAuth,
} from "@trieoh/identityx-sdk-ts-react";

function Header() {
  const { isAuthenticated, auth } = useAuth();

  return isAuthenticated ? (
    <button onClick={() => auth.logout()}>Logout</button>
  ) : null;
}

function LoginPage() {
  return <ModernAuth />;
}
```

All components require `AuthProvider`. The package also exports password
recovery, verification and logout components.

## Server-managed sessions

For server runtimes such as TanStack Start, implement an
`AuthProviderAdapter` that restores the session and creates the auth service
through server functions. Keep IdentityX tokens server-side and prefer an
opaque `HttpOnly`, `Secure`, `SameSite` session cookie.
