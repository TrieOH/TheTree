import { createFileRoute } from '@tanstack/react-router'
import { allPublicEventsQueryOptions } from '@/features/events/api'

export const Route = createFileRoute('/events/')({
  loader: async ({ context }) => {
    return await context.queryClient.ensureQueryData(allPublicEventsQueryOptions())
  },
})