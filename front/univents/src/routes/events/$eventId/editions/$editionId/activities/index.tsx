import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auths/lib/route-guard'

export const Route = createFileRoute(
  '/events/$eventId/editions/$editionId/activities/',
)({
  beforeLoad: requireAuth,
})
