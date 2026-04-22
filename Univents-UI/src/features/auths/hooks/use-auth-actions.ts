import { useNavigate } from "@tanstack/react-router"
import { useAuth } from "@soramux/identityx-sdk-ts/react"
import { toast } from "sonner"
import { getContext } from "@/integrations/tanstack-query/root-provider"

export function useAuthActions() {
  const { auth, isAuthenticated } = useAuth()
  const navigate = useNavigate()

  const handleLoginSuccess = async (redirect?: string) => {
    await navigate({ to: redirect ?? "/", replace: true })
    toast.success("Login successful!")
  }

  const handleLogout = async () => {
    const { queryClient } = getContext()

    queryClient.clear()

    const response = await auth.logout()

    if (!response.success) {
      toast.error("Logout failed")
      return
    }

    await navigate({ to: "/", replace: true })
    toast.success("Logout successful!")
  }

  return {
    handleLoginSuccess,
    handleLogout,
    isAuthenticated,
  }
}