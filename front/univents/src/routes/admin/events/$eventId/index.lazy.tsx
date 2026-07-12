import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute('/admin/events/$eventId/')({
  component: EventOverviewRoute,
})

function EventOverviewRoute() {
  return <div className="min-h-[50vh]" />
}
