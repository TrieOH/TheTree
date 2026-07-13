import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { CertificationTemplateEditor } from '@/features/certifications/ui/CertificationTemplateEditor'

export const Route = createFileRoute('/admin/events/$eventId_/editions/$editionId/certifications/editor')({
  beforeLoad: requireAuth,
  component: RouteComponent,
})

function RouteComponent() {
  const { eventId, editionId } = Route.useParams()
  return <CertificationTemplateEditor eventId={eventId} editionId={editionId} />
}
