import { requireAuth } from '#/features/auths/lib/route-guard'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/')({
  beforeLoad: (ctx) => {
    requireAuth(ctx);
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/"!</div>
}
