import NamespaceLayout from '#/features/admin/ui/namespace-layout'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$namespaceID')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <NamespaceLayout>
      <Outlet />
    </NamespaceLayout>
  )
}
