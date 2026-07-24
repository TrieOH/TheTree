import { createFileRoute, redirect } from '@tanstack/react-router'
import { publicEventBySlugQueryOptions } from '@/features/events/api'
import { activeEditionQueryOptions, pastEditionsQueryOptions, upcomingEditionsQueryOptions } from '@/features/editions/api'

export const Route = createFileRoute('/events/$slug/')({
  loader: async ({ context, params }) => {
    const event = await context.queryClient.ensureQueryData(
      publicEventBySlugQueryOptions(params.slug),
    )
    if (!event) throw redirect({ to: '/events' })

    void context.queryClient.prefetchQuery(activeEditionQueryOptions(event.id))
    void context.queryClient.prefetchQuery(upcomingEditionsQueryOptions(event.id))
    void context.queryClient.prefetchQuery(pastEditionsQueryOptions(event.id))
    return event
  }
})
