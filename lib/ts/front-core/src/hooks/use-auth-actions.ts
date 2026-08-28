import { useNavigate, useRouter } from "@tanstack/react-router"
import { useQueryClient } from "@tanstack/react-query"
import { useAuth } from "@trieoh/identityx-sdk-ts/react"
import { toast } from "sonner"

/**
 * Hook that provides standard login/logout actions.
 *
 * Usage:
 *   const { handleLoginSuccess, handleLogout } = useAuthActions()
 */
export function useAuthActions() {
  const { auth: authManager } = useAuth()
  const router = useRouter()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const handleAuthAction = async (
    isAuthenticated: boolean,
    destination: string,
    successMessage: string,
    performAction?: () => Promise<{ success: boolean }>,
  ) => {
    const auth = router.options.context.auth
    if (!auth) {
      toast.error("Auth Initialization Failed")
      return
    }

    if (performAction) {
      const response = await performAction()
      if (!response.success) {
        toast.error("Auth action failed")
        return
      }
    }

    router.update({
      context: {
        ...router.options.context,
        auth: { ...auth, isAuthenticated },
      },
    })

    await navigate({ to: destination, replace: true })
    toast.success(successMessage)
  }

  const handleLoginSuccess = async (redirect?: string) => {
    await handleAuthAction(true, redirect || "/admin", "Login successful!")
    queryClient.invalidateQueries()
  }

  const handleLogoutTo = async (destination: string) => {
    await handleAuthAction(
      false,
      destination,
      "Logout successful!",
      () => authManager.logout(),
    )
    queryClient.clear()
  }

  const handleLogout = () => handleLogoutTo("/")

  return { handleLoginSuccess, handleLogout, handleLogoutTo }
}
