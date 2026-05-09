import { createFileRoute } from '@tanstack/react-router'
import { RelationshipsFeature } from '#/features/relationships/ui/relationships-feature'

export const Route = createFileRoute('/$envId/relationships')({
  component: RouteComponent,
})

function RouteComponent() {
  const { envId } = Route.useParams()
  return <RelationshipsFeature envId={envId} />
}
