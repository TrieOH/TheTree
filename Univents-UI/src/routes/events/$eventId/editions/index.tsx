import { createFileRoute } from '@tanstack/react-router'
import { allEditionsQueryOptions } from '@/features/editions/api'

export const Route = createFileRoute('/events/$eventId/editions/')({
  loader: async ({ context: ctx, params }) => {
    return await ctx.queryClient.ensureQueryData(
      allEditionsQueryOptions(params.eventId)
    )
  },
})
