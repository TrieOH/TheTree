import { allNamespaceFormsQueryOptions, createFormOnNamespaceFn } from '#/features/forms/api'
import { formCreateSchema } from '#/features/forms/model'
import type { FormCreateI, FormI } from '#/features/forms/model'
import { FormList } from '#/features/forms/ui/form-list'
import FormModal from '#/widgets/modal/form-modal'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/$namespaceID/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { namespaceID } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const queryClient = useQueryClient()
  const userId = auth?.auth.profile()?.id || ''

  const { data: forms = [], isLoading } = useQuery(
    allNamespaceFormsQueryOptions(namespaceID, userId)
  )

  const { mutate: createForm, isPending: isPendingCreate } = useMutation({
    mutationFn: (data: FormCreateI) => createFormOnNamespaceFn(data, namespaceID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allNamespaceFormsQueryOptions(namespaceID, userId).queryKey,
          (old: FormI[] = []) => [...old, response.data],
        )
        setIsCreateOpen(false)
        toast.success('Form created successfully')
      }
    },
    onError: (error: Error) => toast.error(error.message)
  })

  if (isLoading) {
    return (
      <div className="space-y-8 animate-in fade-in duration-500">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div className="space-y-1">
            <div className="h-9 w-48 bg-muted animate-pulse rounded-none" />
            <div className="h-5 w-64 bg-muted animate-pulse rounded-none" />
          </div>
          <div className="h-10 w-full sm:w-36 bg-muted animate-pulse rounded-none" />
        </div>

        <div className="grid gap-6 grid-cols-[repeat(auto-fill,minmax(min(100%,320px),1fr))]">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-40 border-2 border-border/50 bg-card/50 animate-pulse rounded-none" />
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-700">
      <FormList
        forms={forms}
        openModal={() => setIsCreateOpen(true)}
        namespaceID={namespaceID}
      />
      <FormModal<FormCreateI>
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
      />
    </div>
  )
}
