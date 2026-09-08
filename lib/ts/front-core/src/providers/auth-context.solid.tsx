import { useRouter } from "@tanstack/solid-router"
import { useAuth, type AuthService } from "@trieoh/identityx-sdk-ts-solid"
import { createEffect, type ParentProps } from "solid-js"

export type RouterAuth = {
  auth: AuthService
  isAuthenticated: boolean
}

/** Keeps the TanStack Router context in sync with IdentityX auth state. */
export function AuthContextUpdater(props: ParentProps) {
  const auth = useAuth()
  const router = useRouter()

  createEffect(
    () => {
      const context = router.options.context as { auth?: RouterAuth }
      return [auth.auth, auth.isAuthenticated(), context.auth, context] as const
    },
    ([authService, isAuthenticated, currentRouterAuth, context]) => {
      const currentAuth = {
        auth: authService,
        isAuthenticated,
      }

      if (
        currentRouterAuth?.auth !== currentAuth.auth ||
        currentRouterAuth?.isAuthenticated !== currentAuth.isAuthenticated
      ) {
        router.update({
          context: {
            ...context,
            auth: currentAuth,
          },
        })

        if (currentRouterAuth?.isAuthenticated !== currentAuth.isAuthenticated) {
          setTimeout(() => void router.invalidate(), 0)
        }
      }
    },
  )

  return props.children
}
