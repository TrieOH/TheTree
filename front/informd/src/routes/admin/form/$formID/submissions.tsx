import { useMemo, useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  allNamespacesArchivedFormsQueryOptions,
  allNamespacesFormsQueryOptions,
} from '#/features/namespaces/api'
import {
  allUserArchivedFormsQueryOptions,
  allUserFormsQueryOptions,
} from '#/features/forms/api'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import FormAdminHeader from '#/features/forms/ui/form-admin-header'
import type { FormI } from '#/features/forms/model'
import FormHeader from '#/features/submissions/ui/form-header'
import { SubmissionDetail } from '#/features/submissions/ui/submission-details'
import { deriveSubmissions } from '#/features/submissions/model'
import type { SubmissionSummaryI } from '#/features/submissions/model';
import ResponseCard from '#/features/submissions/ui/response-card'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { allFormsResponsesQueryOptions } from '#/features/submissions/api'
import { FileText } from 'lucide-react'

export const Route = createFileRoute('/admin/form/$formID/submissions')({
  component: SubmissionsComponent,
})

function SubmissionsComponent() {
  const { formID } = Route.useParams()
  const { namespaceID } = Route.useSearch()
  const queryClient = useQueryClient()

  const [selectedResponder, setSelectedResponder] = useState<string | null>(null);
  const [filterText, setFilterText] = useState("");

  const { data: fullForm = null, isLoading: isFullFormLoading } = useQuery(allFormsResponsesQueryOptions(formID, namespaceID))

  const submissions = useMemo(() => fullForm ? deriveSubmissions(fullForm) : [], [fullForm]);

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

  const responseCount = submissions.length

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

  if (isFullFormLoading || formsQuery.isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-100 gap-4">
        <div className="size-10 rounded-full border-4 border-primary/20 border-t-primary animate-spin" />
        <p className="text-sm font-medium text-muted-foreground animate-pulse">Loading submissions...</p>
      </div>
    )
  }

  if (!fullForm || !form) {
    return (
      <div className="flex flex-col items-center justify-center min-h-100 gap-4 p-8 text-center">
        <div className="size-16 rounded-2xl bg-muted flex items-center justify-center">
          <FileText className="size-8 text-muted-foreground" />
        </div>
        <div className="space-y-1">
          <h2 className="text-lg font-bold text-foreground">Form not found</h2>
          <p className="text-sm text-muted-foreground max-w-62.5">
            We couldn't find the form you're looking for or it might have been deleted.
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className='space-y-2 animate-in fade-in duration-500'>
      <FormHeader form={fullForm.form} />
      <PaginatedContainer<SubmissionSummaryI>
        items={filteredSubmissions}
        layout="list"
        pageSize={10}
        sortFields={[
          { key: "responder", label: "Respondent" },
          { key: "completed_at", label: "Date" },
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
              steps={fullForm.steps}
              isSelected={selectedResponder === item.responder}
              onClick={() => setSelectedResponder(item.responder)}
            />
          ))
        }
      />
      <SubmissionDetail
        fullForm={fullForm}
        responder={selectedResponder}
        onClose={() => setSelectedResponder(null)}
      />
    </div>
  )
}
