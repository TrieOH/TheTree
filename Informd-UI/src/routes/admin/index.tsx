import { ProjectList } from '#/features/projects/ui/project-list'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <ProjectList 
      openModal={() => {}}
      projects={[
        {id: "23232", name: "", created_at: "", owner_id: "", scope_id: "", updated_at: ""}
      ]}
    />
  )
}
