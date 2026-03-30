import { requireAuth } from '@/features/auth/lib/route-guard';
import FieldEditor from '@/features/schema-version/ui/FieldEditor'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/projects/$projectId/schemas/$schemaId/editor/')({
  beforeLoad: requireAuth,
  component: RouteComponent,
  staticData: {components: {header: "schemas/editor"}}
})

function RouteComponent() {
  return <FieldEditor />
}
