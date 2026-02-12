import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/schemas/editor/')({
  component: RouteComponent,
  staticData: {components: {header: "schemas/editor"}}
})

function RouteComponent() {
  return <div>Hello "/schemas/editor/"!</div>
}
