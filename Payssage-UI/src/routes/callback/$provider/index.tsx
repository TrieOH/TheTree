import { getProviderCallbackFn } from '#/features/oauth/api'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { useEffect, useRef } from 'react'
import { toast } from 'sonner'
import z from 'zod'

const queryParams = z.object({
  code: z.string().optional(),
  state: z.string().optional(),
})

export const Route = createFileRoute('/callback/$provider/')({
  validateSearch: (search) => queryParams.parse(search),
  component: CallbackComponent,
})

function CallbackComponent() {
  const { provider } = Route.useParams()
  const { code, state } = Route.useSearch()
  const navigate = Route.useNavigate()

  const executed = useRef(false)

  useEffect(() => {
    if (executed.current) return
    executed.current = true

    if (!code || !state) { navigate({ to: '/' }); return }

    getProviderCallbackFn(code, state, provider).then(res => {
      if (res.success) {
        const callbackURL = new URL(res.data.url)
        const redirect_url = callbackURL.searchParams.get('redirect_url')!

        const target = new URL(redirect_url)
        callbackURL.searchParams.forEach((value, key) => {
          if (key !== 'redirect_url') target.searchParams.set(key, value)
        })
        window.location.href = target.toString()
      } else {
        toast.error('Failed to connect provider.')
        navigate({ to: '/' })
      }
    }).catch(() => {
      toast.error('Failed to connect provider.')
      navigate({ to: '/' })
    })
  }, [code, state, provider, navigate])

  return <RoutePendingComponent />
}

function RoutePendingComponent() {
  const { provider } = Route.useParams()
  return (
    <div className="flex flex-col items-center justify-center min-h-100 space-y-4 animate-in fade-in duration-500">
      <Loader2 className="w-10 h-10 text-primary animate-spin" />
      <div className="text-center space-y-1">
        <h3 className="text-xl font-black uppercase tracking-tighter">Connecting {provider}</h3>
        <p className="text-muted-foreground text-[10px] font-black uppercase tracking-widest animate-pulse">Please wait while we sync your account...</p>
      </div>
    </div>
  )
}
