import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/$project/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div></div>
}
