import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/transactions')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/transactions"!</div>
}
