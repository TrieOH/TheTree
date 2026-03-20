import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/events/$eventId/editions/$editionId/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/events/$eventId/editions/$editionId/"!</div>
}
