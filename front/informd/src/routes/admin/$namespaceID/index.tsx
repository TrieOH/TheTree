import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'

export const Route = createFileRoute('/admin/$namespaceID/')({
  component: RouteComponent,
})

function RouteComponent() {
  // const { namespaceID } = Route.useParams()

  useLayoutHeader(
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Forms</h1>
        <p className="text-sm text-muted-foreground">
          Manage your namespace forms
        </p>
      </div>
      <button className="inline-flex items-center gap-2 rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground shadow-sm hover:bg-primary/90 transition-colors">
        <Plus className="size-4" />
        New Form
      </button>
    </div>,
  )

  return (
    <main className="animate-in fade-in slide-in-from-bottom-4 duration-700">

      dwdw
      {/* <FormModal<FormCreateI>
        title="Create Form"
        description="Give your form a title to identify it."
        buttonTitle="Create Form"
        schema={formCreateSchema}
        formId="create-form-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createForm}
        fields={[
          {
            name: 'title',
            label: 'Form Title',
            type: 'text',
            placeholder: 'e.g. Contact Form',
          },
        ]}
        disabled={isPendingCreate}
      /> */}
    </main>
  )
}
