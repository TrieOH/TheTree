# IdentityX SDK — TypeScript

Framework-agnostic client for the TrieOH IdentityX API. It contains the API
client, authentication services, token storage and shared types. It does not
include React components or hooks.

## Installation

```bash
npm install @trieoh/identityx-sdk-ts
```

## Configuration

```ts
import { configure } from "@trieoh/identityx-sdk-ts";

configure({
  PROJECT_ID: "your-project-id",
  BASE_URL: "https://api.trieoh.com",
});
```

The SDK also reads `VITE_TRIEOH_AUTH_PROJECT_ID`,
`NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID` and `PUBLIC_TRIEOH_AUTH_PROJECT_ID` when
available.

## Authentication service

```ts
import { Api, createAuthService } from "@trieoh/identityx-sdk-ts";

const auth = createAuthService(new Api());
const response = await auth.login("user@example.com", "password");

if (response.success) console.log("Authenticated");
```

The package exports the shared `Api`, `AuthService`, token utilities, profile
types and authentication types. For React bindings, install the separate
`@trieoh/identityx-sdk-ts-react` package.
