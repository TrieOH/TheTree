import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/admin/events/$eventId/editions/$editionId/',
)({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/events/$eventId/editions/$editionId/"!</div>
}
