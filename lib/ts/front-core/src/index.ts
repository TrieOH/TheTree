// Providers
export { PostHogProvider } from "./providers/posthog"
export type { PostHogConfig } from "./providers/posthog"
export {
  TanStackQueryProvider,
  createQueryClient,
} from "./providers/tanstack-query"
export type { QueryClientConfig } from "./providers/tanstack-query"
export { AuthContextUpdater } from "./providers/auth-context"

// Route guards
export { requireAuth, requireGuest, requireSetup, requireSetupNotDone } from "./guards/route-guards"

// Hooks
export { useAuthActions } from "./hooks/use-auth-actions"
