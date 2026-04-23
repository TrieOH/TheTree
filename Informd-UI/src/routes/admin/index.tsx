import { requireAuth } from '#/features/auths/lib/route-guard'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
  beforeLoad: requireAuth
})

function RouteComponent() {
  return <div>Hello "/admin/"!</div>
}
