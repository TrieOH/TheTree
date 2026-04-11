import { useRouter } from "@tanstack/react-router"
import { useAuth } from "@soramux/node-auth-sdk/react"
import { useLayoutEffect } from "react"
import type { ReactNode } from "react"

export function AuthContextUpdater({ children }: { children: ReactNode }) {
  const auth = useAuth()
  const router = useRouter()

  useLayoutEffect(() => {
    const currentRouterAuth = router.options.context.auth

    if (currentRouterAuth !== auth) {
      router.update({
        context: {
          ...router.options.context,
          auth: auth
        }
      })

      if (currentRouterAuth?.isAuthenticated !== auth.isAuthenticated)
        void router.invalidate()
    }
  }, [auth, router])

  return <>{children}</>
}
