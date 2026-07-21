import { createFileRoute } from '@tanstack/react-router'
import { TransactionsDashboard } from '#/features/payment-intents/ui/transactions-dashboard'
import { listAllIntentsQueryOptions } from '#/features/payment-intents/api'
import { useQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/admin/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  const { data: intents = [] } = useQuery(listAllIntentsQueryOptions())
  return <div className='p-6'><TransactionsDashboard intents={intents} /></div>
}
