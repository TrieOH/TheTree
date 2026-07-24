import { useMemo } from 'react'
import { useMultiStepForm } from '@/widgets/multi-step-form/hooks/use-multi-step-form'
import { MultiStepFormModal } from '@/widgets/multi-step-form/ui/multi-step-form-modal'
import { createEditionFormSteps } from '../model/edition-form-steps'
import { editionCreateSchema } from '../model'
import type {
  EditionCreateInputI,
  EditionCreateOutputI,
  EditionI,
} from '../model'

export interface ManageEditionModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreate: (
    values: EditionCreateOutputI,
  ) => Promise<EditionI | null | boolean> | EditionI | null | boolean
}

const emptyDefaultValues: EditionCreateInputI = {
  name: '',
  slug: '',
  starts_at: '',
  ends_at: '',
}

export function ManageEditionModal({
  open,
  onOpenChange,
  onCreate,
}: ManageEditionModalProps) {
  const steps = useMemo(() => createEditionFormSteps(), [])
  const controller = useMultiStepForm({
    schema: editionCreateSchema,
    steps,
    defaultValues: emptyDefaultValues,
    resetOnSuccessValues: emptyDefaultValues,
    onSubmit: async (values) => Boolean(await onCreate(values)),
    onSubmitSuccess: () => onOpenChange(false),
  })

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Criar Nova Edição"
      controller={controller}
      submitLabel="Criar Edição"
    />
  )
}
