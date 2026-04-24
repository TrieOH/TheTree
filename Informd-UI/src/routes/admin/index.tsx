import { allProjectsQueryOptions, createProjectFn } from '#/features/projects/api';
import { projectCreateSchema   } from '#/features/projects/model'
import type {ProjectCreateI, ProjectI} from '#/features/projects/model';
import { ProjectList } from '#/features/projects/ui/project-list'
import FormModal from '#/widgets/modal/form-modal'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/')({
  component: RouteComponent,
})

function RouteComponent() {
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const queryClient = useQueryClient();

  const { data: projects = [], isLoading } = useQuery(allProjectsQueryOptions())

  const { mutate: createProject, isPending: isPendingCreate } = useMutation({
    mutationFn: createProjectFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allProjectsQueryOptions().queryKey,
          (old: ProjectI[] = []) => [...old, response.data],
        )
        setIsCreateOpen(false)
        toast.success('Project created successfully')
      }
    },
  })

  if (isLoading) {
    return (
      <div className="space-y-8 animate-in fade-in duration-500">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div className="space-y-1">
            <div className="h-9 w-48 bg-muted animate-pulse rounded-sm" />
            <div className="h-5 w-64 bg-muted animate-pulse rounded-sm" />
          </div>
          <div className="h-10 w-full sm:w-36 bg-muted animate-pulse rounded-sm" />
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-35 border border-border bg-card p-5 space-y-4">
              <div className="space-y-2">
                <div className="h-6 w-3/4 bg-muted animate-pulse rounded-sm" />
                <div className="h-4 w-1/4 bg-muted animate-pulse rounded-sm" />
              </div>
              <div className="flex justify-between items-center mt-auto">
                <div className="h-3 w-20 bg-muted animate-pulse rounded-sm" />
                <div className="h-3 w-24 bg-muted animate-pulse rounded-sm" />
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-700">
      <ProjectList 
        openModal={() => setIsCreateOpen(true)}
        projects={projects}
      />
      <FormModal<ProjectCreateI>
        title="Create Project"
        description="Give your project a name to identify it."
        buttonTitle="Create Project"
        schema={projectCreateSchema}
        formId="create-project-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createProject}
        fields={[
          {
            name: "name",
            label: "e.g. My Team Project",
            type: "text",
          }
        ]}
        disabled={isPendingCreate}
      />
    </div>
  )
}
