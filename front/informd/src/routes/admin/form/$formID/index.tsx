import { allFormsStepsQueryOptions, createStepFn } from '#/features/steps/api'
import { stepCreateSchema } from '#/features/steps/model'
import type { StepCreateI, StepI } from '#/features/steps/model';
import { StepCard } from '#/features/steps/ui/step-card'
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
import { PaginatedContainer } from '#/widgets/pagination/paginated-container-grid'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/admin/form/$formID/')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const { formID } = Route.useParams()
  const { namespaceID } = Route.useSearch()
  const { data: steps = [] } = useQuery(allFormsStepsQueryOptions(formID, namespaceID))

  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [filter, setFilter] = useState('')

  const count = steps.length
  const header = useMemo(() => (
    <div className="flex items-start justify-between">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">Steps</h1>
        <p className="text-sm text-muted-foreground">
          {count === 0
            ? 'No steps yet in this form'
            : `${count} step${count !== 1 ? 's' : ''} in this form`}
        </p>
      </div>
    </div>
  ), [count])
  useLayoutHeader(header)

  const filteredSteps = steps.filter((step) => {
    const search = filter.toLowerCase().trim()

    if (!search) return true

    return (
      step.title.toLowerCase().includes(search) ||
      (step.description?.toLowerCase().includes(search) ?? false)
    )
  })

  const { mutate: addStepToForm, isPending: isCreating } = useMutation({
    mutationFn: (data: StepCreateI) => createStepFn(data, formID, namespaceID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allFormsStepsQueryOptions(formID, namespaceID).queryKey,
          (oldData: StepI[] = []) => [...oldData, response.data]
        )
        setIsCreateOpen(false)
        toast.success(response.message || "Step added successfully")
      } else toast.error(response.message || "Failed to add step")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div>
      <PaginatedContainer<StepI>
        items={filteredSteps}
        layout='list'
        pageSize={10}
        sortFields={[
          { key: "title", label: "Title" },
          { key: "position_hint", label: "Position" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by description or title…"
        itemLabel="steps"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            size="icon"
            variant="outline"
            className="sm:w-auto px-3 rounded-sm"
          >
            <Plus size={16} />
            <span className="hidden sm:inline ml-2">Add Step</span>
          </Button>
        }
        renderItems={(slice) => slice.map(item => {
          return <StepCard key={item.id} step={item} />
        })}
      />

      <FormModal<StepCreateI>
        title="Add Step"
        description="Create a new step for this form."
        buttonTitle="Add Step"
        schema={stepCreateSchema}
        formId="add-step-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={addStepToForm}
        fields={[
          {
            name: 'title',
            label: 'Step Title',
            type: 'text',
            placeholder: 'e.g. Personal Information',
          },
          {
            name: 'position_hint',
            label: 'Position Hint',
            type: 'number',
            min: 1,
            placeholder: 'e.g. 1 (This determines the order of steps in the form)',
          },
          {
            name: 'description',
            label: 'Step Description',
            type: 'textarea',
            rows: 4,
            placeholder: 'e.g. Collect basic personal information from the user.',
          }
        ]}
        disabled={isCreating}
      />
    </div>
  )
}
