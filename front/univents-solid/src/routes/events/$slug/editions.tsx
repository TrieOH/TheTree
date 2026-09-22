import { createFileRoute } from '@tanstack/solid-router'

export const Route = createFileRoute('/events/$slug/editions')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/events/$slug/editions"!</div>
}
