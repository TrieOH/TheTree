import type { FormCreateI, FormI } from '#/features/forms/model';
import { FormsView } from '#/features/forms/ui/forms-view'
import { allUserArchivedFormsQueryOptions, allUserFormsQueryOptions, createFormFn } from '#/features/forms/api'
import { allNamespacesArchivedFormsQueryOptions, allNamespacesFormsQueryOptions, createFormOnNamespaceFn } from '#/features/namespaces/api'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { toast } from 'sonner'
import z from 'zod'

const formsSearchSchema = z.object({
  namespaceID: z.string().optional(),
})

export const Route = createFileRoute('/admin/form/')({
  validateSearch: (search) => formsSearchSchema.parse(search),
  component: RouteComponent,
})

function RouteComponent() {
  const { namespaceID } = Route.useSearch()
  const queryClient = useQueryClient()

  // Select query and mutation based on namespaceID
  const queryOptions = namespaceID
    ? allNamespacesFormsQueryOptions(namespaceID)
    : allUserFormsQueryOptions()
  const archivedQueryOptions = namespaceID
    ? allNamespacesArchivedFormsQueryOptions(namespaceID)
    : allUserArchivedFormsQueryOptions()

  const { data: forms = [] } = useQuery(queryOptions)
  const { data: archivedForms = [] } = useQuery(archivedQueryOptions)
  const allForms = [...forms, ...archivedForms]

  const { mutate: createForm, isPending: isCreating } = useMutation({
    mutationFn: (data: FormCreateI) =>
      namespaceID
        ? createFormOnNamespaceFn(namespaceID, data)
        : createFormFn(data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(queryOptions.queryKey, (oldData: FormI[] = []) => {
          return [response.data, ...oldData];
        })
        queryClient.setQueryData(archivedQueryOptions.queryKey, (oldData: FormI[] = []) => {
          return oldData;
        })
        toast.success(response.message || "Form created successfully")
      } else toast.error(response.message || "Failed to create form")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div className="p-6">
      <FormsView
        forms={allForms}
        onCreate={createForm}
        isCreating={isCreating}
        title={namespaceID ? "Namespace Forms" : "My Forms"}
        description={namespaceID ? "in this namespace" : "associated with your account"}
      />
    </div>
  )
}
