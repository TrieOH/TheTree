import { ProjectLayout } from '#/features/admin/ui/project-layout'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$project')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <ProjectLayout>
      <Outlet />
    </ProjectLayout>
  )
}
