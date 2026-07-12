import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute('/admin/events/$eventId/editions/')({
  component: EditionsRoute,
})

function EditionsRoute() {
  return <div className="min-h-[50vh]" />
}
