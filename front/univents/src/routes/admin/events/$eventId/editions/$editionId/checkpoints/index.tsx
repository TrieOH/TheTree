import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { allCheckpointsQueryOptions } from '@/features/checkpoints/api'
import { CheckpointsTab } from '@/features/checkpoints/ui/CheckpointsTab'
import { AdminShell, PageHeader } from '@/shared/ui'
import { editionQueryOptions } from '@/features/editions/api'
import { Skeleton } from '@/shared/ui/shadcn/skeleton'

export const Route = createFileRoute(
  '/admin/events/$eventId/editions/$editionId/checkpoints/',
)({
  component: CheckpointsPage,
})

function CheckpointsPage() {
  const { eventId, editionId } = Route.useParams()

  const { data: edition, isLoading: isLoadingEdition } = useQuery(
    editionQueryOptions(eventId, editionId)
  )

  const { data: checkpoints = [], isLoading: isLoadingCheckpoints } = useQuery(
    allCheckpointsQueryOptions(eventId, editionId)
  )

  const breadcrumbs = (
    <div className="flex items-center gap-2 text-sm text-gray-500">
      <Link to="/admin/events" className="hover:text-gray-900 transition-colors">Eventos</Link>
      <span className="text-gray-300">/</span>
      <Link
        to="/admin/events/$eventId/editions"
        params={{ eventId }}
        className="hover:text-gray-900 transition-colors"
      >
        Edições
      </Link>
      <span className="text-gray-300">/</span>
      <span className="text-gray-900">Checkpoints</span>
    </div>
  )

  if (isLoadingEdition || isLoadingCheckpoints) {
    return (
      <AdminShell breadcrumbs={breadcrumbs}>
        <div className="space-y-6">
          <Skeleton className="h-10 w-48" />
          <div className="space-y-4">
            <Skeleton className="h-24 w-full" />
            <Skeleton className="h-24 w-full" />
            <Skeleton className="h-24 w-full" />
          </div>
        </div>
      </AdminShell>
    )
  }

  return (
    <AdminShell breadcrumbs={breadcrumbs}>
      <PageHeader
        title={`Checkpoints - ${edition?.edition_name}`}
        subtitle="Gerencie os pontos de controle desta edição"
      />

      <div className="mt-6">
        <CheckpointsTab checkpoints={checkpoints} />
      </div>
    </AdminShell>
  )
}
