import { useRouter } from "@tanstack/react-router"
import { useAuth } from "@soramux/identityx-sdk-ts/react"
import { useLayoutEffect } from "react"

export function AuthSynchronizer({ children }: { children: React.ReactNode }) {
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

      if (currentRouterAuth?.isAuthenticated !== auth.isAuthenticated) router.invalidate()
    }
  }, [auth, router])

  return <>{children}</>
}
