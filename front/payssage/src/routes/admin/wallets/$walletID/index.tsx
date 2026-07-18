import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { walletByIdQueryOptions } from '#/features/wallets/api'
import { WalletCollector } from '#/features/wallets/ui/wallet-collector'
import { CreateIntentForm } from '#/features/payment-intents/ui/create-intent-form'

export const Route = createFileRoute('/admin/wallets/$walletID/')({
  component: CollectorPage,
})

function CollectorPage() {
  const { walletID } = Route.useParams()
  const { data: wallet } = useQuery(walletByIdQueryOptions(walletID))
  return wallet ? (
    <div className="space-y-6">
      <WalletCollector wallet={wallet} />
      <CreateIntentForm walletId={walletID} />
    </div>
  ) : null
}
