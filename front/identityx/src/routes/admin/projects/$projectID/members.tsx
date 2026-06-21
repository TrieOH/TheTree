import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/projects/$projectID/members')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/projects/$projectID/members"!</div>
}
