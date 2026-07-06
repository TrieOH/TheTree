import { createFileRoute, Link, useRouter } from '@tanstack/react-router'
import { ModernAuth } from '@trieoh/identityx-sdk-ts/react'
import z from 'zod';
import { requireGuest } from '@/features/auth/lib/route-guard';
import { toast } from 'sonner';
import { ArrowLeft } from 'lucide-react';
import { cn } from '@/shared/lib/utils';

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(''),
})

export const Route = createFileRoute('/auth/')({
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: RouteComponent,
})

function RouteComponent() {

  const navigate = Route.useNavigate()
  const router = useRouter()
  const search = Route.useSearch()

  const handleLoginSuccess = async (message?: string) => {
    const auth = router.options.context.auth
    if (auth) {
      router.update({
        context: {
          ...router.options.context,
          auth: { ...auth, isAuthenticated: true },
        },
      })
      const destination = search.redirect || '/admin'
      await navigate({ to: destination, replace: true })
      toast.success(message ?? "Login successful!")
      router.options.context.queryClient.invalidateQueries();
    } else toast.error("Auth Initialization Failed")
  }

  const handleSignUpSuccess = async (message?: string) => {
    toast.success(message ?? "Account successfully created!")
  }

  const handleFailure = async (message: string, trace?: string[]) => {
    const traceMsg = trace?.join("\n").replaceAll("trace: ", "")
    toast.warning(`Auth Failed: ${message}`, { description: traceMsg })
  }

  return (
    <ModernAuth
      initialView='signin'
      onLoginSuccess={handleLoginSuccess}
      onSignUpSuccess={handleSignUpSuccess}
      onFailed={handleFailure}
      providers={["google", "github"]}
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