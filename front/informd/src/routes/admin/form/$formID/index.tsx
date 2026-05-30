import { allFormsStepsQueryOptions, bulkEditStepsFn, createStepFn } from '#/features/steps/api'
import { stepCreateSchema, stepUpdateSchema } from '#/features/steps/model'
import type { StepCreateI, StepI, StepUpdateI } from '#/features/steps/model';
import { StepCarousel } from '#/features/steps/ui/step-carousel';
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import FormModal from '#/widgets/modal/form-modal'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
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
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [editingStep, setEditingStep] = useState<StepI | null>(null)
  const [defaultValues, setDefaultValues] = useState<Partial<StepCreateI>>({})
  const [addContext, setAddContext] = useState<string | null>(null)
  const [focusedStepId, setFocusedStepId] = useState<string | null>(null)
  const [focusKey, setFocusKey] = useState(0)

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

  const openEditModal = (step: StepI) => {
    setEditingStep(step)
    setIsEditOpen(true)
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

  const { mutate: editStep, isPending: isEditing } = useMutation({
    mutationFn: (data: StepUpdateI) => bulkEditStepsFn([data], formID, namespaceID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({ queryKey: allFormsStepsQueryOptions(formID, namespaceID).queryKey })
        setIsEditOpen(false)
        setEditingStep(null)
        toast.success(response.message || "Step updated successfully")
      } else toast.error(response.message || "Failed to update step")
    },
    onError: (error: Error) => toast.error(error.message)
  })

  const { mutate: moveStep } = useMutation({
    mutationFn: ({ stepId, direction }: { stepId: string, direction: 'left' | 'right' }) => {
      const currentStep = steps.find(s => s.id === stepId)
      if (!currentStep) throw new Error("Step not found")

      const neighborPosition = direction === 'left'
        ? currentStep.position_hint - 1
        : currentStep.position_hint + 1

      const neighborStep = steps.find(s => s.position_hint === neighborPosition)
      if (!neighborStep) throw new Error("No adjacent step found")

      const updatedSteps: StepUpdateI[] = [
        { id: currentStep.id, title: currentStep.title, description: currentStep.description, position_hint: neighborStep.position_hint },
        { id: neighborStep.id, title: neighborStep.title, description: neighborStep.description, position_hint: currentStep.position_hint },
      ]

      return bulkEditStepsFn(updatedSteps, formID, namespaceID)
    },
    onSuccess: (response, variables) => {
      if (response.success) {
        queryClient.setQueryData(
          allFormsStepsQueryOptions(formID, namespaceID).queryKey,
          (oldData: StepI[] = []) => {
            const currentStep = oldData.find(s => s.id === variables.stepId)
            if (!currentStep) return oldData

            const neighborPosition = variables.direction === 'left'
              ? currentStep.position_hint - 1
              : currentStep.position_hint + 1

            return oldData.map(s => {
              if (s.id === variables.stepId) {
                return { ...s, position_hint: neighborPosition }
              }
              if (s.position_hint === neighborPosition) {
                return { ...s, position_hint: currentStep.position_hint }
              }
              return s
            }).sort((a, b) => a.position_hint - b.position_hint)
          }
        )
        setFocusedStepId(variables.stepId)
        setFocusKey(k => k + 1)
      } else toast.error(response.message || "Failed to reorder step")
    },
    onError: (error: Error) => toast.error(error.message),
  })

  return (
    <div>
      <StepCarousel
        steps={steps}
        focusedStepId={focusedStepId}
        focusKey={focusKey}
        onMoveStep={(stepId, direction) => moveStep({ stepId, direction })}
        onEditStep={openEditModal}
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

      <FormModal<StepUpdateI>
        title="Edit Step"
        description="Update the step's information."
        buttonTitle="Save Changes"
        schema={stepUpdateSchema}
        formId="edit-step-form"
        isOpen={isEditOpen}
        defaultValues={editingStep || undefined}
        onClose={() => setIsEditOpen(false)}
        onSubmit={editStep}
        fields={[
          {
            name: 'title',
            label: 'Step Title',
            type: 'text',
            placeholder: 'e.g. Personal Information',
          },
          {
            name: 'description',
            label: 'Step Description',
            type: 'textarea',
            rows: 4,
            placeholder: 'e.g. Collect basic personal information from the user.',
          }
        ]}
        disabled={isEditing}
      />
    </div>
  )
}
