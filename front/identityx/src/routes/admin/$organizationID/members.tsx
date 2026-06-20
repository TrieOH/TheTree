import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$organizationID/members')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/$organizationID/members"!</div>
}
