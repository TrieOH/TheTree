import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute(
  '/events/$eventId/editions/$editionId/activities/',
)({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/events/$eventId/editions/$editionId/activities/"!</div>
}
