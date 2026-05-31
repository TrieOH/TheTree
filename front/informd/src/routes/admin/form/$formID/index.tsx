import { allFormsStepsQueryOptions, bulkEditStepsFn, createStepFn } from '#/features/steps/api'
import { stepCreateSchema, stepUpdateSchema } from '#/features/steps/model'
import type { StepCreateI, StepI, StepUpdateI } from '#/features/steps/model';
import { StepCarousel } from '#/features/steps/ui/step-carousel';
import {
  allStepsFieldsQueryOptions,
  createFieldFn,
  bulkEditFieldsFn,
  deleteFieldFn,
} from '#/features/fields/api'
import { createFieldRequestSchema, fieldUpdateRequestSchema } from '#/features/fields/model'
import type {
  CreateFieldRequestI,
  FieldI,
  FieldUpdateI,
} from '#/features/fields/model';
import { useLayoutHeader } from '#/shared/lib/hooks/layout-context'
import FormModal from '#/widgets/modal/form-modal'
import { ConfirmModal } from '#/widgets/modal/modal'
import { useMutation, useQuery, useQueries, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useMemo, useState, useCallback } from 'react'
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

  const [isFieldCreateOpen, setIsFieldCreateOpen] = useState(false)
  const [isFieldEditOpen, setIsFieldEditOpen] = useState(false)
  const [isFieldDeleteOpen, setIsFieldDeleteOpen] = useState(false)
  const [fieldStepContext, setFieldStepContext] = useState<StepI | null>(null)
  const [editingField, setEditingField] = useState<FieldI | null>(null)
  const [deletingField, setDeletingField] = useState<FieldI | null>(null)

  const count = steps.length
  const maxPosition = useMemo(() => {
    if (count === 0) return 0
    return Math.max(...steps.map(s => s.position_hint))
  }, [steps, count])

  // Fetch fields for ALL steps so any step the carousel shows has its fields ready
  const fieldQueries = useQueries({
    queries: steps.map(step => ({
      ...allStepsFieldsQueryOptions(formID, step.id, namespaceID),
      enabled: steps.length > 0,
    })),
  })
  const fieldsByStepId = useMemo<Record<string, FieldI[]>>(() => {
    const map: Record<string, FieldI[]> = {}
    steps.forEach((step, i) => {
      map[step.id] = fieldQueries[i]?.data ?? []
    })
    return map
  }, [steps, fieldQueries])

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

  // ── Field mutations ──────────────────────────────────────────────────────────

  const { mutate: addField, isPending: isFieldCreating } = useMutation({
    mutationFn: ({ data, step }: { data: CreateFieldRequestI; step: StepI }) =>
      createFieldFn(data, formID, step.id, namespaceID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allStepsFieldsQueryOptions(formID, fieldStepContext?.id ?? '', namespaceID).queryKey,
        })
        setIsFieldCreateOpen(false)
        setFieldStepContext(null)
        toast.success(response.message || "Field added successfully")
      } else toast.error(response.message || "Failed to add field")
    },
    onError: (error: Error) => toast.error(error.message),
  })

  const { mutate: editField, isPending: isFieldEditing } = useMutation({
    mutationFn: ({ data, step }: { data: FieldUpdateI; step: StepI }) =>
      bulkEditFieldsFn([data], formID, step.id, namespaceID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allStepsFieldsQueryOptions(formID, editingField?.step_id ?? '', namespaceID).queryKey,
        })
        setIsFieldEditOpen(false)
        setEditingField(null)
        toast.success(response.message || "Field updated successfully")
      } else toast.error(response.message || "Failed to update field")
    },
    onError: (error: Error) => toast.error(error.message),
  })

  const { mutate: deleteFieldMutation, isPending: isFieldDeleting } = useMutation({
    mutationFn: ({ fieldId, stepId }: { fieldId: string; stepId: string }) =>
      deleteFieldFn(fieldId, formID, stepId, namespaceID),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({
          queryKey: allStepsFieldsQueryOptions(formID, deletingField?.step_id ?? '', namespaceID).queryKey,
        })
        setIsFieldDeleteOpen(false)
        setDeletingField(null)
        toast.success(response.message || "Field deleted successfully")
      } else toast.error(response.message || "Failed to delete field")
    },
    onError: (error: Error) => toast.error(error.message),
  })

  // ── Field handlers ───────────────────────────────────────────────────────────

  const openAddFieldModal = useCallback((step: StepI) => {
    setFieldStepContext(step)
    setIsFieldCreateOpen(true)
  }, [])

  const openEditFieldModal = useCallback((field: FieldI) => {
    setEditingField(field)
    setIsFieldEditOpen(true)
  }, [])

  const openDeleteFieldModal = useCallback((field: FieldI) => {
    setDeletingField(field)
    setIsFieldDeleteOpen(true)
  }, [])

  const handleAddFieldSubmit = useCallback(
    (data: CreateFieldRequestI) => {
      if (!fieldStepContext) return
      addField({ data, step: fieldStepContext })
    },
    [addField, fieldStepContext]
  )

  const handleEditFieldSubmit = useCallback(
    (data: FieldUpdateI) => {
      if (!editingField) return
      // Find the step that owns this field
      const step = steps.find(s => s.id === editingField.step_id)
      if (!step) return
      editField({ data, step })
    },
    [editField, editingField, steps]
  )

  const handleDeleteFieldConfirm = useCallback(() => {
    if (!deletingField) return
    deleteFieldMutation({
      fieldId: deletingField.id,
      stepId: deletingField.step_id,
    })
  }, [deleteFieldMutation, deletingField])

  return (
    <div>
      <StepCarousel
        steps={steps}
        focusedStepId={focusedStepId}
        focusKey={focusKey}
        fieldsByStepId={fieldsByStepId}
        onMoveStep={(stepId, direction) => moveStep({ stepId, direction })}
        onEditStep={openEditModal}
        onAddAfter={(hint) => openAddModal(hint, `after "${steps.find(s => s.position_hint === hint - 1)?.title || 'step'}"`)}
        onAddField={openAddFieldModal}
        onEditField={openEditFieldModal}
        onDeleteField={openDeleteFieldModal}
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

      {/* Field: Create */}
      <FormModal<CreateFieldRequestI>
        title="Add Field"
        description={
          fieldStepContext
            ? `Add a new field to "${fieldStepContext.title}".`
            : 'Add a new field to this step.'
        }
        buttonTitle="Add Field"
        schema={createFieldRequestSchema}
        formId="add-field-form"
        isOpen={isFieldCreateOpen}
        defaultValues={{
          position_hint: (fieldStepContext
            ? fieldsByStepId[fieldStepContext.id].length + 1
            : 1
          ),
          required: false,
        }}
        onClose={() => { setIsFieldCreateOpen(false); setFieldStepContext(null) }}
        onSubmit={handleAddFieldSubmit}
        fields={[
          {
            name: 'key',
            label: 'Field Key',
            type: 'text',
            placeholder: 'e.g. full_name',
          },
          {
            name: 'title',
            label: 'Field Label',
            type: 'text',
            placeholder: 'e.g. Full Name',
          },
          {
            name: 'description',
            label: 'Description',
            type: 'textarea',
            rows: 3,
            placeholder: 'e.g. Enter your full name as it appears on your ID.',
          },
          {
            name: 'type',
            label: 'Field Type',
            type: 'select',
            placeholder: 'Select a type…',
            options: [
              { label: 'String', value: 'string' },
              { label: 'Email', value: 'email' },
              { label: 'Integer', value: 'int' },
              { label: 'Float', value: 'float' },
              { label: 'Boolean', value: 'bool' },
              { label: 'Date', value: 'date' },
              { label: 'Time', value: 'time' },
              { label: 'Datetime', value: 'datetime' },
              { label: 'Select', value: 'select' },
              { label: 'File', value: 'file' },
              { label: 'Phone', value: 'phone' },
              { label: 'URL', value: 'url' },
            ],
          },
          {
            name: 'required',
            label: 'Required',
            type: 'select',
            placeholder: 'No',
            options: [
              { label: 'Yes', value: 'true' },
              { label: 'No', value: 'false' },
            ],
          },
          {
            name: 'placeholder',
            label: 'Placeholder',
            type: 'text',
            placeholder: 'e.g. Type your answer here…',
          },
          {
            name: 'default_value',
            label: 'Default Value',
            type: 'text',
            placeholder: 'e.g. John Doe',
          },
          {
            name: 'position_hint',
            label: 'Position (auto)',
            type: 'number',
            disabled: true,
          },
        ]}
        disabled={isFieldCreating}
      />

      {/* Field: Edit */}
      <FormModal<FieldUpdateI>
        title="Edit Field"
        description={editingField ? `Update "${editingField.title}".` : 'Update the field.'}
        buttonTitle="Save Changes"
        schema={fieldUpdateRequestSchema}
        formId="edit-field-form"
        isOpen={isFieldEditOpen}
        defaultValues={editingField || undefined}
        onClose={() => { setIsFieldEditOpen(false); setEditingField(null) }}
        onSubmit={handleEditFieldSubmit}
        fields={[
          {
            name: 'key',
            label: 'Field Key',
            type: 'text',
            placeholder: 'e.g. full_name',
          },
          {
            name: 'title',
            label: 'Field Label',
            type: 'text',
            placeholder: 'e.g. Full Name',
          },
          {
            name: 'description',
            label: 'Description',
            type: 'textarea',
            rows: 3,
            placeholder: 'e.g. Enter your full name as it appears on your ID.',
          },
          {
            name: 'type',
            label: 'Field Type',
            type: 'select',
            placeholder: 'Select a type…',
            options: [
              { label: 'String', value: 'string' },
              { label: 'Email', value: 'email' },
              { label: 'Integer', value: 'int' },
              { label: 'Float', value: 'float' },
              { label: 'Boolean', value: 'bool' },
              { label: 'Date', value: 'date' },
              { label: 'Time', value: 'time' },
              { label: 'Datetime', value: 'datetime' },
              { label: 'Select', value: 'select' },
              { label: 'File', value: 'file' },
              { label: 'Phone', value: 'phone' },
              { label: 'URL', value: 'url' },
            ],
          },
          {
            name: 'required',
            label: 'Required',
            type: 'select',
            placeholder: 'No',
            options: [
              { label: 'Yes', value: 'true' },
              { label: 'No', value: 'false' },
            ],
          },
          {
            name: 'placeholder',
            label: 'Placeholder',
            type: 'text',
            placeholder: 'e.g. Type your answer here…',
          },
          {
            name: 'default_value',
            label: 'Default Value',
            type: 'text',
            placeholder: 'e.g. John Doe',
          },
        ]}
        disabled={isFieldEditing}
      />

      {/* Field: Delete confirmation */}
      <ConfirmModal
        isOpen={isFieldDeleteOpen}
        onClose={() => { setIsFieldDeleteOpen(false); setDeletingField(null) }}
        onConfirm={handleDeleteFieldConfirm}
        title="Delete Field"
        description={
          deletingField
            ? `Are you sure you want to delete "${deletingField.title}"? This action cannot be undone.`
            : 'Are you sure you want to delete this field?'
        }
        confirmText="Delete"
        variant="destructive"
        isLoading={isFieldDeleting}
      />
    </div>
  )
}
