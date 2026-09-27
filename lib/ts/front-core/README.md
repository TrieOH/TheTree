# @trieoh/front-core

Framework-free core shared by the Tria apps: the IdentityX BFF (cookie session,
transport, operation handlers) and the OpenTelemetry browser/server wiring.

Nothing here imports React or Solid, so the same code runs in the React apps'
TanStack Start server functions, in the Solid app's Worker, in a plain Node
script and in tests. The framework bindings sit in siblings:

| Package | Contains |
|---|---|
| `@trieoh/front-core` (this) | `auth/bff/*`, `tracing/*` |
| [`@trieoh/front-core-react`](../front-core-react) | TanStack Start client adapter + BFF server binding, providers, guards, `useAuthActions` |
| [`@trieoh/front-core-solid`](../front-core-solid) | `QueryClient` provider + auth-context updater |

Most consumers import a subpath rather than the root barrel:

```ts
import { createBffHandler } from "@trieoh/front-core/auth/bff/handler";
import { withSpan } from "@trieoh/front-core/tracing/browser";
```

The name suffix is the framework contract — see *Package boundaries* in the root
`README.md`, enforced by `pnpm check:boundaries`.
