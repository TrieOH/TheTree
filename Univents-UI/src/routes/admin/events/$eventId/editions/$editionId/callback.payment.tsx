import { createFileRoute } from '@tanstack/react-router'
import z from 'zod'
import { toast } from 'sonner'
import { useQueryClient } from '@tanstack/react-query'
import { useEffect, useRef } from 'react'
import type { EditionI } from '@/features/editions/model'
import {
  allAdminEditionsQueryOptions,
  connectPaymentAccountToEditionFn
} from '@/features/editions/api'
import WaveSpinnerLoading from '@/shared/ui/loader/WaveSpinnerLoading'

const queryParams = z.object({
  credential_id: z.string().optional(),
  provider: z.string().optional(),
  public_key: z.string().optional()
})

export const Route = createFileRoute(
  '/admin/events/$eventId/editions/$editionId/callback/payment',
)({
  validateSearch: (search) => queryParams.parse(search),
  component: PaymentCallbackComponent,
})

function PaymentCallbackComponent() {
  const { credential_id, provider, public_key } = Route.useSearch()
  const { eventId, editionId } = Route.useParams()
  const navigate = Route.useNavigate()
  const queryClient = useQueryClient()
  const hasCalled = useRef(false)

  useEffect(() => {
    if (hasCalled.current) return
    hasCalled.current = true

    if (!credential_id || !provider || !public_key) {
      toast.error('Missing required payment credentials')
      void navigate({ to: '/admin/events/$eventId/editions', params: { eventId } })
      return
    }

    connectPaymentAccountToEditionFn(
      eventId, editionId, credential_id, provider, public_key
    ).then(res => {
      if (res.success) {
        queryClient.setQueryData(
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
        toast.success('Connected payment account to edition with success.')
      } else toast.error('Failed to connect payment account to edition.')
    }).catch(() => toast.error('An error occurred during payment account connection.'))
      .finally(() => navigate({ to: '/admin/events/$eventId/editions', params: { eventId } }))
  }, [credential_id, provider, public_key, eventId, editionId, navigate, queryClient])

  return (
    <div className="flex flex-col items-center justify-center min-h-100 space-y-4 animate-in fade-in duration-500">
      <div className="text-center space-y-4">
        <WaveSpinnerLoading text='Connecting...' />
        <p className="text-muted-foreground text-[10px] font-black tracking-widest animate-pulse">Please wait while we sync your account...</p>
      </div>
    </div>
  )
}
