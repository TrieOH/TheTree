import { createFileRoute } from '@tanstack/react-router'
import { SignatureEditor } from '@/features/signatures/ui/SignatureEditor'

export const Route = createFileRoute('/admin/events/$eventId_/editions/$editionId/signatures/editor')({
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  return <SignatureEditor eventId={eventId} editionId={editionId} />
}
