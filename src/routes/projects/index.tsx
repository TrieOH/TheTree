import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auth/lib/route-guard';
import { ProjectDialog } from '@/features/project/ui/ProjectDialog';
import ProjectCard from '@/features/project/ui/ProjectCard';
import { projectsQueryOptions } from '@/features/project/api';
import { useSuspenseQuery } from '@tanstack/react-query';

export const Route = createFileRoute('/projects/')({
  beforeLoad: requireAuth,
  staticData: {
    components: {
      header: "projects"
    }
  },
  loader: async ({ context: { queryClient }}) => {
    await queryClient.ensureQueryData(projectsQueryOptions)
    return {}
  },
  component: RouteComponent,
})

function RouteComponent() {
  const projectsQuery = useSuspenseQuery(projectsQueryOptions)
  const projects = projectsQuery.data
  
  return (
    <main className="w-full bg-background flex flex-col items-center my-4">
      <div className="text-center space-y-1 mb-7">
        <h1 className="font-bold text-3xl">Your Projects</h1>
        <p className="font-extralight text-sm">
          Manage your projects configurations
        </p>
      </div>
      <div className="max-w-7xl grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {projects.map((project) => (
          <ProjectCard key={project.id} data={project} />
        ))}
      </div>
      <ProjectDialog />
    </main>
  )
}
