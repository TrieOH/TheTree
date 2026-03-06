# NodeAuth

SDK for integrating with the TrieOH authentication ecosystem.

## Installation

```bash
npm install @trieoh/node-auth-sdk
# or
yarn add @trieoh/node-auth-sdk
# or
bun add @trieoh/node-auth-sdk
```

## Configuration (Vite / React)

To use the SDK in a React project (Vite, Next.js, or CRA), wrap your application with `AuthProvider`.

### Option 1: Environment Variables (Recommended)

The SDK automatically looks for these variables:

- `VITE_TRIEOH_AUTH_PROJECT_ID` (Vite)
- `NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID` (Next.js)
- `PUBLIC_TRIEOH_AUTH_PROJECT_ID` (General)

```tsx
import { AuthProvider } from '@trieoh/node-auth-sdk/react';

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
import { configure } from '@trieoh/node-auth-sdk';

configure({
  PROJECT_ID: 'your-id',
  BASE_URL: 'https://your-api.com'
});
```

## Components

The SDK provides ready-to-use components:

```tsx
import { SignIn, SignUp } from '@trieoh/node-auth-sdk/react';

// Example usage
const LoginPage = () => <SignIn />;
const RegisterPage = () => <SignUp />;
```

## Hooks

You can access the authentication state anywhere in your application:

```tsx
import { useAuth } from '@trieoh/node-auth-sdk/react';

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

## Server-Side

For server-side operations (Node.js / Next.js API routes), use `createServerAuth`:

```ts
import { createServerAuth } from '@trieoh/node-auth-sdk/server';

const auth = createServerAuth();
// auth.assignRoleByNameToUser(...)
```

---