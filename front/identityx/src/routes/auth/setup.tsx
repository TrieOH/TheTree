import { createFileRoute, Link, useNavigate, useRouter } from '@tanstack/react-router'
import { ModernSetup } from '@trieoh/identityx-sdk-ts/react'
import { requireGuest } from '@/features/auth/lib/route-guard';
import { toast } from 'sonner';
import { ArrowLeft } from 'lucide-react';
import { cn } from '@/shared/lib/utils';

export const Route = createFileRoute('/auth/setup')({
  beforeLoad: requireGuest,
  component: RouteComponent,
})

function RouteComponent() {
  const navigate = useNavigate()
  const router = useRouter()

  const handleSetupSuccess = async (message?: string) => {
    toast.success(message ?? "Initial setup completed successfully!")
    const auth = router.options.context.auth
    if (auth) {
      router.update({
        context: {
          ...router.options.context,
          auth: { ...auth, isAuthenticated: true },
        },
      })
      await navigate({ to: '/projects', replace: true })
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
      backLink={
        <Link
          to="/"
          className={cn(
            "inline-flex items-center gap-0.5 text-sm text-muted-foreground",
            "hover:text-foreground transition-colors"
          )}
        >
          <ArrowLeft size={20} /> Back to Home
        </Link>
      }
    />
  )
}