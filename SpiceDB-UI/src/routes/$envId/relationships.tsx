import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/$envId/relationships')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <main className="flex flex-col h-(--content-height) min-w-75 border-l">

    </main>
  )
}
