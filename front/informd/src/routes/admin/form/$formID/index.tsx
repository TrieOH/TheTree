import { allFormsStepsQueryOptions, createStepFn } from '#/features/steps/api'
import { stepCreateSchema } from '#/features/steps/model'
import type { StepCreateI, StepI } from '#/features/steps/model';
import { StepCarousel } from '#/features/steps/ui/step-carousel';
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import { Button } from '#/shared/ui/shadcn/button'
import FormModal from '#/widgets/modal/form-modal'
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
  const [defaultValues, setDefaultValues] = useState<Partial<StepCreateI>>({})
  const [addContext, setAddContext] = useState<string | null>(null)

  const count = steps.length
  const maxPosition = useMemo(() => {
    if (count === 0) return 0
    return Math.max(...steps.map(s => s.position_hint))
  }, [steps, count])

  const openAddModal = (requestedHint: number, contextName?: string) => {
    // Check if hint is already taken
    const isTaken = steps.some(s => s.position_hint === requestedHint)
    const finalHint = isTaken ? maxPosition + 1 : requestedHint

    setDefaultValues({ position_hint: finalHint })
    setAddContext(contextName || null)
    setIsCreateOpen(true)
  }

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
      <Button
        size="sm"
        className="gap-2"
        onClick={() => openAddModal(maxPosition + 1)}
      >
        <Plus className="w-4 h-4" />
        Add Step
      </Button>
    </div>
  ), [count, maxPosition])
  useLayoutHeader(header)

  const { mutate: addStepToForm, isPending: isCreating } = useMutation({
    mutationFn: (data: StepCreateI) => createStepFn(data, formID, namespaceID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allFormsStepsQueryOptions(formID, namespaceID).queryKey,
          (oldData: StepI[] = []) => {
            const newData = [...oldData, response.data]
            return newData.sort((a, b) => a.position_hint - b.position_hint)
          }
        )
        setIsCreateOpen(false)
        toast.success(response.message || "Step added successfully")
      } else toast.error(response.message || "Failed to add step")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  return (
    <div>
      <StepCarousel
        steps={steps}
        onAddBefore={(hint) => openAddModal(hint, `before "${steps.find(s => s.position_hint === hint + 1)?.title || 'step'}"`)}
        onAddAfter={(hint) => openAddModal(hint, `after "${steps.find(s => s.position_hint === hint - 1)?.title || 'step'}"`)}
      />

      <FormModal<StepCreateI>
        title="Add Step"
        description={addContext ? `This step will be created ${addContext}.` : "Create a new step for this form."}
        buttonTitle="Add Step"
        schema={stepCreateSchema}
        formId="add-step-form"
        isOpen={isCreateOpen}
        defaultValues={defaultValues}
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
            label: 'Position (System Managed)',
            type: 'number',
            disabled: true,
            placeholder: 'Position is automatically assigned',
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
