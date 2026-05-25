import { FormCard } from '#/features/forms/ui/form-card'
import { allNamespacesFormsQueryOptions, createFormOnNamespaceFn } from '#/features/namespaces/api'
import { formCreateOnNamespaceSchema, type FormCreateOnNamespaceI, type FormI } from '#/features/namespaces/model'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'


export const Route = createFileRoute('/admin/$namespaceID/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { namespaceID } = Route.useParams()
  const queryClient = useQueryClient()
  const { data: forms = [] } = useQuery(allNamespacesFormsQueryOptions(namespaceID))

  const [filter, setFilter] = useState('')
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  const count = forms.length

  const header = useMemo(() => (
    <div className="flex items-start justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Forms</h1>
        <p className="text-sm text-muted-foreground">
          {count === 0
            ? 'No forms yet in this namespace'
            : `${count} form${count !== 1 ? 's' : ''} in this namespace`}
        </p>
      </div>
    </div>
  ), [count])

  useLayoutHeader(header)

  const filteredForms = forms.filter((form) => {
    const search = filter.toLowerCase().trim()

    if (!search) return true

    return (
      form.title.toLowerCase().includes(search) ||
      form.created_at.toLowerCase().includes(search) ||
      form.updated_at.toLowerCase().includes(search) ||
      form.status.toLowerCase().includes(search)
    )
  })

  const { mutate: createForm, isPending: isCreating } = useMutation({
    mutationFn: (data: FormCreateOnNamespaceI) => createFormOnNamespaceFn(namespaceID, data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(allNamespacesFormsQueryOptions(namespaceID).queryKey, (oldData: FormI[] = []) => {
          return [response.data, ...oldData];
        })
        setIsCreateOpen(false)
        toast.success(response.message || "Form created successfully")
      } else toast.error(response.message || "Failed to create form")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div>
      <PaginatedContainer<FormI>
        items={filteredForms}
        className='w-full'
        layout='flex'
        pageSize={10}
        sortFields={[
          { key: "title", label: "Title" },
          { key: "created_at", label: "Created At" },
          { key: "updated_at", label: "Updated At" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by title…"
        itemLabel="forms"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            size="icon"
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
          >
            <Plus size={16} />
            <span className="hidden sm:inline ml-2">Add Form</span>
          </Button>
        }
        renderItems={(slice) => slice.map(item => {
          return (
            <FormCard key={item.id} data={item} />
          )
        })}
      />
      <FormModal<FormCreateOnNamespaceI>
        title="Create Form"
        description="Give your form a title to identify it."
        buttonTitle="Create Form"
        schema={formCreateOnNamespaceSchema}
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
        disabled={isCreating}
      />
    </div>
  )
}
