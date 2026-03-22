import { createFileRoute } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import z from 'zod'
import { toast } from 'sonner'
import type { EditionI } from '@/features/editions/model'
import {
  allAdminEditionsQueryOptions,
  connectPaymentAccountToEditionFn
} from '@/features/editions/api'


const queryParams = z.object({
  credential_id: z.string().optional(),
  provider: z.string().optional(),
  public_key: z.string().optional()
})

export const Route = createFileRoute(
  '/admin/events/$eventId/editions/$editionId/callback/payment',
)({
  validateSearch: (search) => queryParams.parse(search),
  beforeLoad: async ({ search, params, context }) => {
    const { credential_id, provider, public_key } = search
    const { eventId, editionId } = params
    if (!credential_id || !provider || !public_key)
      throw Route.redirect({
        to: '/admin/events/$eventId/editions/$editionId',
        params: { eventId, editionId }
      })
    const res = await connectPaymentAccountToEditionFn(
      eventId, editionId, credential_id, provider, public_key
    );
    if (res.success) {
      context.queryClient.setQueryData(
        allAdminEditionsQueryOptions(eventId).queryKey,
        (old: EditionI[] = []) =>
          old.map((ed) =>
            ed.id === editionId
              ? {
                ...ed,
                trie_payments_credential_id: credential_id,
                trie_payments_provider: provider,
                trie_payments_provider_public_key: public_key
              }
              : ed
          )
      )
      toast.success('connected payment account to edition with success.')
    } else toast.error('Failed to connect payment account to edition.')
    throw Route.redirect({
      to: '/admin/events/$eventId/editions/$editionId',
      params: { eventId, editionId }
    })
  },
  pendingComponent: RoutePendingComponent,
  component: () => null,
})

function RoutePendingComponent() {
  return (
    <div className="flex flex-col items-center justify-center min-h-100 space-y-4 animate-in fade-in duration-500">
      <Loader2 className="w-10 h-10 text-primary animate-spin" />
      <div className="text-center space-y-1">
        <h3 className="text-xl font-black uppercase tracking-tighter">Connecting</h3>
        <p className="text-muted-foreground text-[10px] font-black uppercase tracking-widest animate-pulse">Please wait while we sync your account...</p>
      </div>
    </div>
  )
}