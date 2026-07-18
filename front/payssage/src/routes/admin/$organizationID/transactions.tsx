import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$organizationID/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/$organizationID/transactions"!</div>
}
