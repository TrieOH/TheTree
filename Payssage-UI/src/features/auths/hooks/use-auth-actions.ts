import { useNavigate, useRouter } from '@tanstack/react-router'
import { useAuth } from '@soramux/node-auth-sdk/react'
import { toast } from 'sonner'
import { getContext } from '@/app/integrations/tanstack-query/root-provider'

export function useAuthActions() {
  const { auth: authManager } = useAuth()
  const router = useRouter()
  const navigate = useNavigate()

  const handleAuthAction = async (
    isAuthenticated: boolean,
    destination: string,
    successMessage: string,
    performAction?: () => Promise<{ success: boolean }>,
  ) => {
    const auth = router.options.context.auth
    if (!auth) {
      toast.error('Auth Initialization Failed')
      return
    }

    if (performAction) {
      const response = await performAction()
      if (!response.success) {
        toast.error('Auth action failed')
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
    await handleAuthAction(true, redirect || '/admin', 'Login successful!')
  }

  const handleLogout = async () => {
    const { queryClient } = getContext()
    queryClient.clear()
    await handleAuthAction(false, '/', 'Logout successful!', () => authManager.logout())
  }

  return { handleLoginSuccess, handleLogout }
}
