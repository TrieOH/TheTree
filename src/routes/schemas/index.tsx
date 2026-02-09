import { navigationStore } from '@/features/navigation';
import { SchemaDialog } from '@/features/schema/ui/SchemaDialog';
import { createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from "@tanstack/react-store"

export const Route = createFileRoute('/schemas/')({
  beforeLoad: async () => {
    const currentProjectId = navigationStore.state.currentProjectId;
    if (typeof window !== 'undefined' && !currentProjectId) throw redirect({ to: '/projects' });
  },
  component: SchemaPage,
  staticData: {components: {header: "schemas"}}
})

function SchemaPage() {
  const currentProjectId = useStore(navigationStore, (state) => state.currentProjectId);
  return (
    <main>
      <SchemaDialog project_id={currentProjectId || ""}/>
    </main>
    
  )
}

