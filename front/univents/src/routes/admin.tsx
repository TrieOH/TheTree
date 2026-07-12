import { createFileRoute, Outlet } from '@tanstack/react-router'
import { requireAuth } from '@/features/auths/lib/route-guard'
import { AdminLayout } from '@/features/admin/ui/admin-layout'

export const Route = createFileRoute('/admin')({
  beforeLoad: requireAuth,
  component: AdminLayoutWrapper,
})

function AdminLayoutWrapper() {
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  )
}
