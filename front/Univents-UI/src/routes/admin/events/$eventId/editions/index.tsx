import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auths/lib/route-guard'

export const Route = createFileRoute('/admin/events/$eventId/editions/')({
  beforeLoad: requireAuth,
})
