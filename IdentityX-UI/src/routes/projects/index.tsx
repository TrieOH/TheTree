import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from '@/features/auth/lib/route-guard';
import { ProjectDialog } from '@/features/project/ui/ProjectDialog';
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
    <main className="w-full bg-background">
      <ProjectDialog />
    </main>
  )
}
