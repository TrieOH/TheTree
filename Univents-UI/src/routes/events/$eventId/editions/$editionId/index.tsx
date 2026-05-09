import { createFileRoute } from '@tanstack/react-router'
import { editionQueryOptions } from '@/features/editions/api'

export const Route = createFileRoute('/events/$eventId/editions/$editionId/')({
  loader: async ({ context: ctx, params }) => {
    return await ctx.queryClient.ensureQueryData(
      editionQueryOptions(params.eventId, params.editionId)
    )
  },
})
