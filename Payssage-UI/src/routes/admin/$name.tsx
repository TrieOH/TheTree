import { WorkspaceLayout } from '#/features/admin/ui/workspace-layout'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$name')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <WorkspaceLayout>
      <Outlet />
    </WorkspaceLayout>
  )
}
