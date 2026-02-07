import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auth/lib/route-guard';
import { ProjectDialog } from '@/features/project/ui/ProjectDialog';
import ProjectCard from '@/features/project/ui/ProjectCard';
// import { createServerFn } from '@tanstack/react-start';

// createServerFn

export const Route = createFileRoute('/projects/')({
  beforeLoad: requireAuth,
  staticData: {
    components: {
      header: "projects"
    }
  },
  component: RouteComponent,
  // loader
})

function RouteComponent() {
  return (
    <main className="w-full bg-background flex flex-col items-center mt-4">
      <div className="text-center space-y-1 mb-7">
        <h1 className="font-bold text-3xl">Your Projects</h1>
        <p className="font-extralight text-sm">
          Manage your projects configurations
        </p>
      </div>
      <div className="max-w-7xl">
        <ProjectCard
          data={{
            id: "0",
            is_active: true,
            project_name: "Univents",
            metadata: {},
            created_at: "2026-02-06T23:06:11.608471Z",
            updated_at: "2026-02-07T23:06:11.608471Z",
            owner_id: ""
          }}
        />
      </div>
      <ProjectDialog />
    </main>
  )
}
