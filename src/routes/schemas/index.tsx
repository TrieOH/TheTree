import { requireAuth } from '@/features/auth/lib/route-guard';
import { navigationStore } from '@/features/navigation';
import { schemasQueryOptions } from '@/features/schema/api';
import { SchemaDialog } from '@/features/schema/ui/SchemaDialog';
import { useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from "@tanstack/react-store"

export const Route = createFileRoute('/schemas/')({
  beforeLoad: async (ctx) => {
    requireAuth(ctx)
    const currentProjectId = navigationStore.state.currentProjectId;
    if (typeof window !== 'undefined' && !currentProjectId) throw redirect({ to: '/projects' });
  },
  loader: async ({ context: { queryClient }}) => {
    const currentProjectId = navigationStore.state.currentProjectId;
    await queryClient.ensureQueryData(schemasQueryOptions(currentProjectId || ""));
    return {}
  },
  component: SchemaPage,
  staticData: {components: {header: "schemas"}}
})

function SchemaPage() {
  const currentProjectId = useStore(navigationStore, (state) => state.currentProjectId);
  const { data: schemas } = useSuspenseQuery(schemasQueryOptions(currentProjectId || ""));
  return (
    <main className="w-full bg-background flex flex-col items-center my-4">
      <div className="text-center space-y-1 mb-7">
        <h1 className="font-bold text-3xl">Schemas</h1>
        <p className="font-extralight text-sm">
          Manage your schemas configurations for this project
        </p>
      </div>
      <div className="max-w-7xl w-full xs:px-4">
        <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'>
          {schemas.map((schema) => (
            <p key={schema.id}>{schema.title}</p>
          ))} 
          {/* {projects.map((project) => (
            <ProjectCard key={project.id} data={project} />
          ))} */}
          {/* <ProjectAddButton
            onCreate={projectActions.openCreate}
          /> */}
        </div>
      </div>
      <SchemaDialog project_id={currentProjectId || ""}/>
    </main>
  )
}

