# @trieoh/front-core-solid

Solid bindings for [`@trieoh/front-core`](../front-core): the `QueryClient`
provider and the auth-context updater the Solid app renders at its root. Nothing
here touches the BFF internals — those live in the neutral package.

```tsx
import {
  AuthContextUpdater,
  TanStackQueryProvider,
  createQueryClient,
} from "@trieoh/front-core-solid";

const queryClient = createQueryClient();

<AuthContextUpdater session={session}>
  <TanStackQueryProvider client={queryClient}>{children}</TanStackQueryProvider>
</AuthContextUpdater>;
```

The React counterpart is [`@trieoh/front-core-react`](../front-core-react); the
name suffix is the framework contract — see *Package boundaries* in the root
`README.md`.
