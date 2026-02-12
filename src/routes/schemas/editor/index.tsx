import FieldEditor from '@/features/schema-version/ui/FieldEditor'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/schemas/editor/')({
  component: RouteComponent,
  staticData: {components: {header: "schemas/editor"}}
})

function RouteComponent() {
  return <FieldEditor />
}
