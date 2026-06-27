import { allProjectsQueryOptions, createProjectFn } from '@/features/project/api'
import type { ProjectCreateI, ProjectI } from '@/features/project/model'
import { ProjectsView } from '@/features/project/ui/ProjectsView'
import { useLayoutHeader } from '@trieoh/ui-base'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useMemo } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/$organizationID/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { organizationID } = Route.useParams()
  const queryClient = useQueryClient()
  const { data: projects = [] } = useQuery(allProjectsQueryOptions(organizationID))

  const count = projects.length

  const header = useMemo(() => (
    <div className="flex items-start justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Projects</h1>
        <p className="text-sm text-muted-foreground">
          {count === 0
            ? 'No projects yet in this organization'
            : `${count} project${count !== 1 ? 's' : ''} in this organization`}
        </p>
      </div>
    </div>
  ), [count])

  useLayoutHeader(header)

  const { mutate: createProject, isPending: isCreating } = useMutation({
    mutationFn: (data: ProjectCreateI) => createProjectFn(data, organizationID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(allProjectsQueryOptions(organizationID).queryKey, (oldData: ProjectI[] = []) => {
          return [response.data, ...oldData];
        })
        toast.success(response.message || "Project created successfully")
      } else toast.error(response.message || "Failed to create project")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <ProjectsView
      projects={projects}
      onCreate={createProject}
      isCreating={isCreating}
      title="" // Title is handled by layout header
      description="" // Description is handled by layout header
    />
  )
}
