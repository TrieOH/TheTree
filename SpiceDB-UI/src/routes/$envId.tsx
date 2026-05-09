import AdminLayout from '#/features/admin/ui/admin-layout'
import { SiteHeader } from '#/shared/ui/site-header'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/$envId')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <SiteHeader />
      <AdminLayout>
        <Outlet />
      </AdminLayout>
    </>
  )
}
