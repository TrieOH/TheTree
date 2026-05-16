import { createFileRoute } from '@tanstack/react-router'
import { CheckFeature } from '#/features/relationships/ui/check-feature'

export const Route = createFileRoute('/$envId/check')({
  component: RouteComponent,
})

function RouteComponent() {
  const { envId } = Route.useParams()
  return <CheckFeature envId={envId} />
}
