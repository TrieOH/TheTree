import { createFileRoute } from '@tanstack/react-router'
import { TransactionsDashboard } from '#/features/payment-intents/ui/transactions-dashboard'

export const Route = createFileRoute('/admin/wallets/$walletID/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  const { walletID } = Route.useParams()
  return <TransactionsDashboard title="Wallet transactions" description="Payment activity for this wallet." walletId={walletID} />
}
