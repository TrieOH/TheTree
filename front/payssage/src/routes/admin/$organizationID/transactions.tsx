import { createFileRoute } from '@tanstack/react-router'
import { TransactionsDashboard } from '#/features/payment-intents/ui/transactions-dashboard'
import { useQuery } from '@tanstack/react-query'
import { listAllByOrgIntentsQueryOptions } from '#/features/payment-intents/api'

export const Route = createFileRoute('/admin/$organizationID/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  const { organizationID } = Route.useParams()

  const { data: intents = [] } = useQuery(listAllByOrgIntentsQueryOptions(organizationID))
  return (
    <TransactionsDashboard
      title="Organization transactions"
      description="Payment activity for wallets in this organization."
      intents={intents}
    />
  )
}
