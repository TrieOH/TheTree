import { allProjectsQueryOptions, createProjectFn } from '@/features/project/api'
import type { ProjectCreateI, ProjectI } from '@/features/project/model'
import { ProjectsView } from '@/features/project/ui/ProjectsView'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { toast } from 'sonner'
import z from 'zod'

const projectsSearchSchema = z.object({
  organizationID: z.string().optional(),
})

export const Route = createFileRoute('/admin/projects/')({
  validateSearch: (search) => projectsSearchSchema.parse(search),
  component: RouteComponent,
})

function RouteComponent() {
  const { organizationID } = Route.useSearch()
  const queryClient = useQueryClient()

  const { data: projects = [] } = useQuery(allProjectsQueryOptions(organizationID))

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
    <div className="p-6">
      <ProjectsView
        projects={projects}
        onCreate={createProject}
        isCreating={isCreating}
        title={organizationID ? "Organization Projects" : "My Projects"}
        description={organizationID ? "in this organization" : "associated with your account"}
      />
    </div>
  )
}
