import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$namespaceID/members')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/$namespaceID/members"!</div>
}
