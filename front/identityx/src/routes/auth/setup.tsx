import { createFileRoute, useNavigate, useRouter } from '@tanstack/react-router'
import { ModernSetup } from '@trieoh/identityx-sdk-ts/react'
import { requireSetupNotDone } from '@/features/auth/lib/route-guard';
import { toast } from 'sonner';

const SETUP_DONE_KEY = 'trieoh_setup_done';

function markSetupComplete() {
  if (typeof window !== "undefined") localStorage.setItem(SETUP_DONE_KEY, "true");
}

export const Route = createFileRoute('/auth/setup')({
  beforeLoad: requireSetupNotDone,
  component: RouteComponent,
})

function RouteComponent() {
  const navigate = useNavigate()
  const router = useRouter()

  const handleSetupSuccess = async (message?: string) => {
    toast.success(message ?? "Initial setup completed successfully!")
    markSetupComplete()
    const auth = router.options.context.auth
    if (auth) {
      router.update({
        context: {
          ...router.options.context,
          auth: { ...auth, isAuthenticated: true },
        },
      })
      await navigate({ to: '/admin', replace: true })
      router.options.context.queryClient.invalidateQueries();
    } else await navigate({ to: '/auth', replace: true })
  }

  const handleFailure = async (message: string, trace?: string[]) => {
    const traceMsg = trace?.join("\n").replaceAll("trace: ", "")
    toast.warning(`Setup Failed: ${message}`, { description: traceMsg })
  }

  return (
    <ModernSetup
      onSuccess={handleSetupSuccess}
      onFailed={handleFailure}
    />
  )
}