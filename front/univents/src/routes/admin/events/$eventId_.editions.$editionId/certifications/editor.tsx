import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/events/$eventId_/editions/$editionId/certifications/editor')({
  component: RouteComponent,
})

function RouteComponent() {
  // const { eventId, editionId } = Route.useParams()
}
