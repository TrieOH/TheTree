import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute(
  '/admin/events/$eventId/editions/$editionId/activities/',
)({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div>Hello "/admin/events/$eventId/editions/$editionId/activities/"!</div>
  )
}
