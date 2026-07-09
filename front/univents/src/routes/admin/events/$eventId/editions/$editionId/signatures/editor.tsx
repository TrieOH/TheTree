import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { SignatureEditor } from '@/features/signatures/ui/SignatureEditor'

export const Route = createFileRoute('/admin/events/$eventId/editions/$editionId/signatures/editor')({
  beforeLoad: requireAuth,
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  return <SignatureEditor eventId={eventId} editionId={editionId} />
}
