import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auth/lib/route-guard';
import { ProjectDialog } from '@/features/project/ui/ProjectDialog';
import ProjectCard from '@/features/project/ui/ProjectCard';
import { projectsQueryOptions } from '@/features/project/api';
import { ProjectsSkeleton } from '@/shared/ui/placeholders/ProjectsSkeleton';
import { ProjectsEmptyState } from '@/features/project/ui/ProjectsEmptyState';
import { projectActions } from '@/features/project/store';
import { cn } from '@/shared/lib/utils';
import { ProjectAddButton } from '@/features/project/ui/ProjectAddButon';
import { useQuery } from '@tanstack/react-query';

export const Route = createFileRoute('/projects/')({
  beforeLoad: requireAuth,
  pendingComponent: ProjectsSkeleton,
  component: RouteComponent,
})

function RouteComponent() {
  const { data: projects = [] } = useQuery(projectsQueryOptions)

  const hasProjects = projects && projects.length > 0

  if (!hasProjects) return (
    <main className={cn(
      "flex justify-center items-center bg-background",
      "w-full h-(--screen--minus-header) px-4"
    )}>
      <ProjectsEmptyState
        onCreate={projectActions.openCreate}
      />
      <ProjectDialog />
    </main>
  )

  return (
    <main className="w-full bg-background flex flex-col items-center my-4">
      <div className="text-center space-y-1 mb-7">
        <h1 className="font-bold text-3xl">Your Projects</h1>
        <p className="font-extralight text-sm">
          Manage your projects configurations
        </p>
      </div>
      <div className="max-w-7xl w-full xs:px-4">
        <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'>
          {projects.map((project) => (
            <ProjectCard key={project.id} data={project} />
          ))}
          <ProjectAddButton
            onCreate={projectActions.openCreate}
          />
        </div>
      </div>
      <ProjectDialog />
    </main>
  )
}
