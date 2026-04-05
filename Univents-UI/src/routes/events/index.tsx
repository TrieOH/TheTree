import { createFileRoute } from '@tanstack/react-router'
import { eventsQueryOptions } from '@/features/events/api'

export const Route = createFileRoute('/events/')({
  loader: async ({ context }) => {
    return await context.queryClient.ensureQueryData(eventsQueryOptions())
  },
})