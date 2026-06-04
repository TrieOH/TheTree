import { useMemo } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  allNamespacesFormsQueryOptions,
  formResponseCountOnNamespaceQueryOptions,
} from '#/features/namespaces/api'
import {
  allUserFormsQueryOptions,
  formResponseCountQueryOptions,
} from '#/features/forms/api'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import FormAdminHeader from '#/features/forms/ui/form-admin-header'
import type { FormI } from '#/features/forms/model'
import { Inbox } from 'lucide-react'

export const Route = createFileRoute('/admin/form/$formID/submissions')({
  component: SubmissionsComponent,
})

function SubmissionsComponent() {
  const { formID } = Route.useParams()
  const { namespaceID } = Route.useSearch()
  const queryClient = useQueryClient()

  const formsQueryKey = useMemo(() =>
    namespaceID
      ? allNamespacesFormsQueryOptions(namespaceID).queryKey
      : allUserFormsQueryOptions().queryKey,
    [namespaceID]
  )

  const formsQuery = useQuery({
    queryKey: formsQueryKey,
    queryFn: namespaceID
      ? allNamespacesFormsQueryOptions(namespaceID).queryFn
      : allUserFormsQueryOptions().queryFn,
  })

  const form = useMemo(() =>
    formsQuery.data?.find((f) => f.id === formID),
    [formsQuery.data, formID]
  )

  const countQuery = useQuery(
    namespaceID
      ? formResponseCountOnNamespaceQueryOptions(namespaceID, formID)
      : formResponseCountQueryOptions(formID)
  )

  const responseCount = countQuery.data?.count ?? 0

  const updateFormData = (updatedForm: FormI) => {
    queryClient.setQueryData(
      formsQueryKey,
      (oldData: FormI[] = []) => oldData.map(f => f.id === formID ? updatedForm : f)
    )
  }

  const header = useMemo(() => {
    if (!form) return null

    return (
      <FormAdminHeader
        title="Submissions"
        description={responseCount === 0 ? `No submissions yet for ${form.title}` : `${responseCount} submission${responseCount !== 1 ? 's' : ''} for ${form.title}`}
        form={form}
        namespaceID={namespaceID}
        responseCount={responseCount}
        onUpdate={updateFormData}
      />
    )
  }, [form, responseCount, namespaceID])

  useLayoutHeader(header)

  if (formsQuery.isLoading) return <div className="p-6 text-sm text-muted-foreground">Loading form...</div>
  if (!form) return <div className="p-6 text-sm text-muted-foreground">Form not found</div>

  return (
    <div className="flex flex-col gap-8">
      {/* Submissions List Placeholder */}
      <div className="border border-dashed rounded-sm p-20 flex flex-col items-center justify-center text-center bg-muted/5">
        <div className="size-12 rounded-full bg-muted/50 flex items-center justify-center mb-4">
          <Inbox className="size-6 text-muted-foreground/40" />
        </div>
        <h2 className="text-base font-semibold">No submissions to display</h2>
        <p className="text-sm text-muted-foreground max-w-xs mt-1">
          Once users start filling out your form, their responses will appear here in a detailed list.
        </p>
      </div>
    </div>
  )
}
