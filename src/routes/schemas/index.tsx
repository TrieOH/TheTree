import { navigationStore } from '@/features/navigation';
import { createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from "@tanstack/react-store"

export const Route = createFileRoute('/schemas/')({
  beforeLoad: async () => {
    const currentProjectId = navigationStore.state.currentProjectId;
    if (typeof window !== 'undefined' && !currentProjectId) throw redirect({ to: '/projects' });
  },
  component: SchemaPage,
  staticData: {components: {header: "projects"}}
})

function SchemaPage() {
  const currentProjectId = useStore(navigationStore, (state) => state.currentProjectId);
  return (
    // Temp
    <div className="p-4">
      <h1 className="text-2xl font-bold">Schema Page</h1>
      <p>This is the schema page for a project.</p>
      {currentProjectId ? (
        <p>Current Project ID: {currentProjectId}</p>
      ) : (
        <p>No project selected.</p>
      )}
    </div>
  )
}

