import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$organizationID/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/$organizationID/"!</div>
}
