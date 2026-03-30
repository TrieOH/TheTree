import { requireAuth } from '#/features/auths/lib/route-guard'
import { createFileRoute, Outlet } from '@tanstack/react-router'
import { AdminLayout } from '#/features/admin/ui/admin-layout'

export const Route = createFileRoute('/admin')({
  beforeLoad: (ctx) => {
    requireAuth(ctx);
  },
  component: AdminLayoutWrapper,
})

function AdminLayoutWrapper() {
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  )
}
