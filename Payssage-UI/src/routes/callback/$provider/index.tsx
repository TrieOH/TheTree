import { getProviderCallbackFn } from '#/features/oauth/api'
import { createFileRoute, redirect, useParams } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import z from 'zod'

const queryParams = z.object({
  code: z.string().optional(),
  state: z.string().optional(),
})

export const Route = createFileRoute('/callback/$provider/')({
  validateSearch: (search) => queryParams.parse(search),
  beforeLoad: async ({ search, params }) => {
    const { code, state } = search
    const { provider } = params
    if (!code || !state) throw redirect({ to: '/' })
    const res = await getProviderCallbackFn(code, state, provider)
    if (res.success) {
      const callbackURL = new URL(res.data.url)
      const redirect_url = callbackURL.searchParams.get('redirect_url')!

      const target = new URL(redirect_url)
      callbackURL.searchParams.forEach((value, key) => {
        if (key !== 'redirect_url') target.searchParams.set(key, value)
      })

      window.location.href = target.toString()
    }
    toast.error('Failed to connect provider.')
    throw redirect({ to: '/' })
  },
  pendingComponent: RoutePendingComponent,
  component: () => null,
})

function RoutePendingComponent() {
  const { provider } = useParams({ from: '/callback/$provider/' })
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
