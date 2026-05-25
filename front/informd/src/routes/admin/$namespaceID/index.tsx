import type { FormCreateI, FormI } from '#/features/forms/model';
import { FormsView } from '#/features/forms/ui/forms-view'
import { allNamespacesFormsQueryOptions, createFormOnNamespaceFn } from '#/features/namespaces/api'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useMemo } from 'react'
import { toast } from 'sonner'


export const Route = createFileRoute('/admin/$namespaceID/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { namespaceID } = Route.useParams()
  const queryClient = useQueryClient()
  const { data: forms = [] } = useQuery(allNamespacesFormsQueryOptions(namespaceID))

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

  const { mutate: createForm, isPending: isCreating } = useMutation({
    mutationFn: (data: FormCreateI) => createFormOnNamespaceFn(namespaceID, data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(allNamespacesFormsQueryOptions(namespaceID).queryKey, (oldData: FormI[] = []) => {
          return [response.data, ...oldData];
        })
        toast.success(response.message || "Form created successfully")
      } else toast.error(response.message || "Failed to create form")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <FormsView
      forms={forms}
      onCreate={createForm}
      isCreating={isCreating}
      title="" // Title is handled by layout header
      description="" // Description is handled by layout header
    />
  )
}
