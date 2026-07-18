import { createFileRoute } from '@tanstack/react-router'
import { TransactionsDashboard } from '#/features/payment-intents/ui/transactions-dashboard'

export const Route = createFileRoute('/admin/$organizationID/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  return <TransactionsDashboard title="Organization transactions" description="Payment activity for wallets in this organization." />
}
