import { useRouter } from "@tanstack/react-router"
import { useAuth } from "@trieoh/identityx-sdk-ts/react"
import { useLayoutEffect } from "react"
import type { ReactNode } from "react"

/**
 * Keeps the TanStack Router context in sync with the IdentityX auth state.
 *
 * Must be rendered inside a `<Router>` and an `<AuthProvider>`.
 *
 * Usage:
 *   <AuthContextUpdater>
 *     {children}
 *   </AuthContextUpdater>
 */
export function AuthContextUpdater({ children }: { children: ReactNode }) {
  const auth = useAuth()
  const router = useRouter()

  useLayoutEffect(() => {
    const currentRouterAuth = router.options.context.auth

    if (currentRouterAuth !== auth) {
      router.update({
        context: {
          ...router.options.context,
          auth,
        },
      })

      if (currentRouterAuth?.isAuthenticated !== auth.isAuthenticated) {
        void router.invalidate()
      }
    }
  }, [auth, router])

  return <>{children}</>
}
