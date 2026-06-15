import { createFileRoute, Link, useNavigate, useRouter, useSearch } from '@tanstack/react-router'
import { ModernAuth } from '@trieoh/identityx-sdk-ts/react'
import z from 'zod';
import { requireGuest } from '@/features/auth/lib/route-guard';
import { toast } from 'sonner';

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(''),
})

export const Route = createFileRoute('/auth/')({
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: RouteComponent,
})

function RouteComponent() {

  const navigate = useNavigate()
  const router = useRouter()
  const search = useSearch({ from: '/auth/' })

  const handleLoginSuccess = async (message?: string) => {
    const auth = router.options.context.auth
    if (auth) {
      router.update({
        context: {
          ...router.options.context,
          auth: { ...auth, isAuthenticated: true },
        },
      })
      const destination = search.redirect || '/projects'
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
    <div className="relative min-h-screen bg-background">
      <Link
        to="/"
        className="absolute top-4 left-4 z-10 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        ← Back to Home
      </Link>
      <ModernAuth
        initialView='signin'
        onLoginSuccess={handleLoginSuccess}
        onSignUpSuccess={handleSignUpSuccess}
        onFailed={handleFailure}
      />
    </div>
  )
}