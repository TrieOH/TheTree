import { requireAuth } from '@/features/auth/lib/route-guard';
import { navigationStore } from '@/features/navigation';
import FieldEditor from '@/features/schema-version/ui/FieldEditor'
import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/schemas/editor/')({
  beforeLoad: async (ctx) => {
    requireAuth(ctx)
    const { currentProjectId, currentSchemaId } = navigationStore.state;
    if (typeof window !== 'undefined') {
      if(!currentProjectId) throw redirect({ to: '/projects' });
      if(!currentSchemaId) throw redirect({ to: '/projects/config' });
    }
  },
  component: RouteComponent,
  staticData: {components: {header: "schemas/editor"}}
})

function RouteComponent() {
  return <FieldEditor />
}
