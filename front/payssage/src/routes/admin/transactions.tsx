import { createFileRoute } from '@tanstack/react-router'
import { TransactionsDashboard } from '#/features/payment-intents/ui/transactions-dashboard'

export const Route = createFileRoute('/admin/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div className='p-6'><TransactionsDashboard /></div>
}
