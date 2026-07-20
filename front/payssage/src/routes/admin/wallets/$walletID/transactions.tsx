import { createFileRoute } from '@tanstack/react-router'
import { TransactionsDashboard } from '#/features/payment-intents/ui/transactions-dashboard'
import { useQuery } from '@tanstack/react-query'
import { listAllByWalletIntentsQueryOptions } from '#/features/payment-intents/api'

export const Route = createFileRoute('/admin/wallets/$walletID/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  const { walletID } = Route.useParams()
  const { data: intents = [] } = useQuery(listAllByWalletIntentsQueryOptions(walletID))
  return (
    <TransactionsDashboard
      title="Wallet transactions"
      description="Payment activity for this wallet."
      intents={intents}
    />
  )
}
