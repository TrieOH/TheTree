import { AdminLayout } from '@/features/admin/ui/AdminLayout'
import { requireAuth } from '@/features/auth/lib/route-guard'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/admin')({
  component: RouteComponent,
  beforeLoad: requireAuth,
})

function RouteComponent() {
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  )
}
