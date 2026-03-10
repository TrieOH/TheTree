import { getProviderCallbackFn } from '#/features/oauth/api'
import { useParams, createFileRoute, redirect } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import z from 'zod'

const queryParams = z.object({
  code: z.string().optional(),
  state: z.string().optional(),
})

export const Route = createFileRoute('/admin/$name/providers/$provider/')({
  pendingComponent: RoutePendingComponent,
  component: () => null,
  validateSearch: (search) => queryParams.parse(search),
  beforeLoad: async ({ search, params }) => {
    const { code, state } = search
    const { name, provider } = params
    if (!code || !state) throw redirect({ to: '/admin/$name/providers', params: { name } })
    const res = await getProviderCallbackFn(code, state, provider)
    if (res.success) {
      const status = new URL(res.data.url).searchParams.get('status') ?? 'failed'

      toast[status === 'success' ? 'success' : 'error'](
        status === 'success'
          ? `Successfully connected ${provider}!`
          : `Failed to connect ${provider}.`
      )

      throw redirect({
        to: '/admin/$name/providers',
        params: { name },
        search: { status, provider },
      })
    }
    toast.error('Failed to connect provider.')
    throw redirect({
      to: '/admin/$name/providers',
      params: { name },
      search: { status: 'failed', provider },
    })
  }
})

function RoutePendingComponent() {
  const { provider } = useParams({ from: '/admin/$name/providers/$provider/' })
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