import { useMultiStepForm } from '@/widgets/multi-step-form/hooks/use-multi-step-form'
import { MultiStepFormModal } from '@/widgets/multi-step-form/ui/multi-step-form-modal'
import {
  eventMemberCreateSchema,
  type EventMemberCreateInput,
  type EventMemberCreateOutput,
} from '../model/member'
import { createEventMemberFormSteps } from '../model/member-form-steps'

interface ManageEventMemberModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreate: (values: EventMemberCreateOutput) => boolean | Promise<boolean>
}

const defaultValues: EventMemberCreateInput = {
  email: '',
  role: 'staff',
}

export function ManageEventMemberModal({
  open,
  onOpenChange,
  onCreate,
}: ManageEventMemberModalProps) {
  const controller = useMultiStepForm({
    schema: eventMemberCreateSchema,
    steps: createEventMemberFormSteps(),
    defaultValues,
    resetOnSuccessValues: defaultValues,
    onSubmit: onCreate,
    onSubmitSuccess: () => onOpenChange(false),
  })

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Adicionar membro"
      controller={controller}
      submitLabel="Adicionar membro"
    />
  )
}
