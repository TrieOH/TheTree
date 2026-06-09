import { useMemo, useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  allNamespacesArchivedFormsQueryOptions,
  allNamespacesFormsQueryOptions,
  formResponseCountOnNamespaceQueryOptions,
} from '#/features/namespaces/api'
import {
  allUserArchivedFormsQueryOptions,
  allUserFormsQueryOptions,
  formResponseCountQueryOptions,
} from '#/features/forms/api'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import FormAdminHeader from '#/features/forms/ui/form-admin-header'
import type { FormI } from '#/features/forms/model'
import { mockForm } from '#/features/submissions/model/mock'
import FormHeader from '#/features/submissions/ui/form-header'
import { SubmissionDetail } from '#/features/submissions/ui/submission-details'
import { deriveSubmissions  } from '#/features/submissions/model'
import type {SubmissionSummaryI} from '#/features/submissions/model';
import ResponseCard from '#/features/submissions/ui/response-card'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'

export const Route = createFileRoute('/admin/form/$formID/submissions')({
  component: SubmissionsComponent,
})

function SubmissionsComponent() {
  const { formID } = Route.useParams()
  const { namespaceID } = Route.useSearch()
  const queryClient = useQueryClient()

  const [selectedResponder, setSelectedResponder] = useState<string | null>(null);
  const [filterText, setFilterText] = useState("");
  const submissions = useMemo(() => deriveSubmissions(mockForm), []);

  const filteredSubmissions = useMemo(() => {
    if (!filterText) return submissions;
    const search = filterText.toLowerCase();
    return submissions.filter(s =>
      s.responder.toLowerCase().includes(search)
    );
  }, [submissions, filterText]);

  const formsQueryKey = useMemo(() =>
    namespaceID
      ? allNamespacesFormsQueryOptions(namespaceID).queryKey
      : allUserFormsQueryOptions().queryKey,
    [namespaceID]
  )
  const archivedFormsQueryKey = useMemo(() =>
    namespaceID
      ? allNamespacesArchivedFormsQueryOptions(namespaceID).queryKey
      : allUserArchivedFormsQueryOptions().queryKey,
    [namespaceID]
  )

  const formsQuery = useQuery({
    queryKey: formsQueryKey,
    queryFn: namespaceID
      ? allNamespacesFormsQueryOptions(namespaceID).queryFn
      : allUserFormsQueryOptions().queryFn,
  })
  const archivedFormsQuery = useQuery({
    queryKey: archivedFormsQueryKey,
    queryFn: namespaceID
      ? allNamespacesArchivedFormsQueryOptions(namespaceID).queryFn
      : allUserArchivedFormsQueryOptions().queryFn,
  })

  const form = useMemo(() =>
    [...(formsQuery.data ?? []), ...(archivedFormsQuery.data ?? [])].find((f) => f.id === formID),
    [formsQuery.data, archivedFormsQuery.data, formID]
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
    queryClient.setQueryData(
      archivedFormsQueryKey,
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
    <div className='space-y-2'>
      <FormHeader form={mockForm.form} />
      <PaginatedContainer<SubmissionSummaryI>
        items={filteredSubmissions}
        layout="list"
        pageSize={10}
        sortFields={[
          { key: "responder", label: "Respondent" },
          { key: "completed_at", label: "Date" },
          { key: "step_id", label: "Step" },
        ]}
        filterValue={filterText}
        onFilterChange={setFilterText}
        filterPlaceholder="Search by respondent email…"
        itemLabel="responses"
        renderItems={(slice) =>
          slice.map((item) => (
            <ResponseCard
              key={item.responder}
              data={item}
              steps={mockForm.steps}
              isSelected={selectedResponder === item.responder}
              onClick={() => setSelectedResponder(item.responder)}
            />
          ))
        }
      />
      <SubmissionDetail
        fullForm={mockForm}
        responder={selectedResponder}
        onClose={() => setSelectedResponder(null)}
      />
    </div>
  )
}
