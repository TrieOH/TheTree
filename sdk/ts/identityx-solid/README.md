# IdentityX SDK — Solid

Solid 2 bindings for the framework-agnostic IdentityX client.

## Installation

```bash
npm install @trieoh/identityx-sdk-ts-solid @trieoh/identityx-sdk-ts solid-js @solidjs/web
```

## Provider

```tsx
import { AuthProvider } from "@trieoh/identityx-sdk-ts-solid";
import "@trieoh/identityx-sdk-ts-solid/styles.css";

export function App() {
  return (
    <AuthProvider>
      <YourRoutes />
    </AuthProvider>
  );
}
```

## Authentication state

```tsx
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";

function Header() {
  const { isAuthenticated, auth } = useAuth();

  return isAuthenticated() ? (
    <button onClick={() => auth.logout()}>Logout</button>
  ) : null;
}
```

The provider accepts `projectId`, `baseURL`, `waitSession` and `fallback`.

## Components

The package exports the Modern components (`ModernAuth`, `ModernSignIn`, `ModernSignUp`,
`ModernForgotPassword`, `ModernResetPassword`, `ModernVerifyEmail`,
`ModernResendVerification`, `ModernSetup`, `ModernProfile` and
`ModernIntrospect`).

```tsx
import { AuthProvider, ModernAuth } from "@trieoh/identityx-sdk-ts-solid";

<AuthProvider projectId="your-project-id">
  <ModernAuth />
</AuthProvider>;
```

The CSS is opt-in so the SDK does not modify global styles without the
application importing `@trieoh/identityx-sdk-ts-solid/styles.css`.
