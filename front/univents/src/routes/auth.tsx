import { createFileRoute, useRouter } from '@tanstack/react-router'
import { ModernAuth } from '@trieoh/identityx-sdk-ts/react'
import z from 'zod'
import { requireGuest } from '@/features/auths/lib/route-guard'

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(''),
})

export const Route = createFileRoute('/auth')({
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: AuthPage,
})

function AuthPage() {
  const navigate = Route.useNavigate()
  const router = useRouter()
  const search = Route.useSearch()

  const handleLoginSuccess = async () => {
    const auth = router.options.context.auth

    if (!auth) return

    router.update({
      context: {
        ...router.options.context,
        auth: { ...auth, isAuthenticated: true },
      },
    })

    const destination = search.redirect || '/profile'
    await navigate({ to: destination, replace: true })
    router.options.context.queryClient.invalidateQueries()
  }

  const handleSignUpSuccess = async () => { }

  const handleFailure = async (message: string) => {
    void message
  }

  return (
    <div className="[&>main]:py-16">
      <ModernAuth
        initialView="signin"
        onLoginSuccess={handleLoginSuccess}
        onSignUpSuccess={handleSignUpSuccess}
        onFailed={handleFailure}
        providers={['google', 'github']}
      />
    </div>
  )
}
