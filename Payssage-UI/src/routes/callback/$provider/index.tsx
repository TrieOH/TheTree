import { getProviderCallbackFn } from '#/features/oauth/api'
import { Button } from '#/shared/ui/shadcn/button'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { useEffect } from 'react'
import z from 'zod'

const queryParams = z.object({
  code: z.string().optional(),
  state: z.string().optional(),
})

type CallbackLoaderData =
  | { ok: true; redirectTo: string }
  | { ok: false; message: string }

export const Route = createFileRoute('/callback/$provider/')({
  validateSearch: (search) => queryParams.parse(search),

  pendingComponent: RoutePendingComponent,
  pendingMs: 200,

  loader: async ({ params, location }): Promise<CallbackLoaderData> => {
    const { provider } = params

    const parsed = queryParams.safeParse(
      Object.fromEntries(new URLSearchParams(location.search))
    )
    if (!parsed.success) return { ok: false, message: 'Invalid OAuth params.' }

    const { code, state } = parsed.data

    if (!code || !state) return { ok: false, message: 'Missing OAuth params.' }

    try {
      const res = await getProviderCallbackFn(code, state, provider)

      if (!res.success) return { ok: false, message: 'Failed to connect provider.' }

      const callbackURL = new URL(res.data.url)
      const redirect_url = callbackURL.searchParams.get('redirect_url')

      if (!redirect_url) return { ok: false, message: 'Missing redirect URL.' }

      const target = new URL(redirect_url)

      callbackURL.searchParams.forEach((value, key) => {
        if (key !== 'redirect_url') target.searchParams.set(key, value)
      })

      return { ok: true, redirectTo: target.toString() }
    } catch {
      return { ok: false, message: 'Unexpected error while connecting provider.' }
    }
  },
  component: CallbackPage,
})

function CallbackPage() {
  const data = Route.useLoaderData()
  const navigate = Route.useNavigate()

  useEffect(() => {
    if (data.ok) window.location.href = data.redirectTo.toString()
  }, [data])

  if (data.ok) {
    return (
      <div className="flex min-h-100 flex-col items-center justify-center space-y-4 animate-in fade-in duration-500">
        <Loader2 className="h-10 w-10 animate-spin text-primary" />
        <div className="text-center space-y-1">
          <h3 className="text-xl font-black uppercase tracking-tighter">
            Connecting provider
          </h3>
          <p className="text-muted-foreground text-[10px] font-black uppercase tracking-widest animate-pulse">
            Redirecting you now...
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-100 flex-col items-center justify-center space-y-4">
      <div className="text-center space-y-2">
        <h3 className="text-xl font-black uppercase tracking-tighter">
          Connection failed
        </h3>
        <p className="text-sm text-muted-foreground">{data.message}</p>
      </div>
      <Button
        className="rounded-md border px-4 py-2 text-sm font-semibold cursor-pointer"
        onClick={() => navigate({ to: "/" })}
      >
        Go home
      </Button>
    </div>
  )
}

function RoutePendingComponent() {
  const { provider } = Route.useParams()

  return (
    <div className="flex flex-col items-center justify-center min-h-100 space-y-4 animate-in fade-in duration-500">
      <Loader2 className="w-10 h-10 text-primary animate-spin" />
      <div className="text-center space-y-1">
        <h3 className="text-xl font-black uppercase tracking-tighter">
          Connecting {provider}
        </h3>
        <p className="text-muted-foreground text-[10px] font-black uppercase tracking-widest animate-pulse">
          Please wait while we sync your account...
        </p>
      </div>
    </div>
  )
}