# NodeAuth

SDK for integrating with the TrieOH authentication ecosystem.

## Installation

```bash
npm install @soramux/identityx-sdk-ts
# or
yarn add @soramux/identityx-sdk-ts
# or
bun add @soramux/identityx-sdk-ts
```

## Configuration (Vite / React)

To use the SDK in a React project (Vite, Next.js, or CRA), wrap your application with `AuthProvider`.

### Option 1: Environment Variables (Recommended)

The SDK automatically looks for these variables:

- `VITE_TRIEOH_AUTH_PROJECT_ID` (Vite)
- `NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID` (Next.js)
- `PUBLIC_TRIEOH_AUTH_PROJECT_ID` (General)

```tsx
import { AuthProvider } from '@soramux/identityx-sdk-ts/react';

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
import { configure } from '@soramux/identityx-sdk-ts';

configure({
  PROJECT_ID: 'your-id',
  BASE_URL: 'https://your-api.com'
});
```

## Components

The SDK provides ready-to-use components:

```tsx
import { SignIn, SignUp } from '@soramux/identityx-sdk-ts/react';

// Example usage
const LoginPage = () => <SignIn />;
const RegisterPage = () => <SignUp />;
```

## Hooks

You can access the authentication state anywhere in your application:

```tsx
import { useAuth } from '@soramux/identityx-sdk-ts/react';

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

---