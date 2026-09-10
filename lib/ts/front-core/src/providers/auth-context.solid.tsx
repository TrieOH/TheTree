import { useRouter } from "@tanstack/solid-router"
import { useAuth, type AuthService } from "@trieoh/identityx-sdk-ts-solid"
import { createEffect, type ParentProps } from "solid-js"

export type RouterSession = {
  service: AuthService
  isAuthenticated: boolean
}

/** Keeps the TanStack Router context in sync with IdentityX auth state. */
export function AuthContextUpdater(props: ParentProps) {
  const auth = useAuth()
  const router = useRouter()

  createEffect(
    () => {
      const context = router.options.context as { session?: RouterSession }
      return [auth.auth, auth.isAuthenticated(), context.session, context] as const
    },
    ([authService, isAuthenticated, currentRouterAuth, context]) => {
      const currentSession = {
        service: authService,
        isAuthenticated,
      }

      if (
        currentRouterAuth?.service !== currentSession.service ||
        currentRouterAuth?.isAuthenticated !== currentSession.isAuthenticated
      ) {
        router.update({
          context: {
            ...context,
            session: currentSession,
          },
        })

        if (currentRouterAuth?.isAuthenticated !== currentSession.isAuthenticated) {
          setTimeout(() => void router.invalidate(), 0)
        }
      }
    },
  )

  return props.children
}
