import { useMemo } from 'react'
import { useMultiStepForm } from '@/widgets/multi-step-form/hooks/use-multi-step-form'
import { MultiStepFormModal } from '@/widgets/multi-step-form/ui/multi-step-form-modal'
import { createEditionFormSteps } from '../model/edition-form-steps'
import { editionCreateSchema } from '../model'
import type { EditionCreateOutputI, EditionCreateInputI, EditionI } from '../model'

export interface ManageEditionModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  edition?: EditionI
  onCreate?: (values: EditionCreateOutputI) => Promise<EditionI | null | boolean> | EditionI | null | boolean
  onUpdate?: (id: string, values: EditionCreateOutputI) => Promise<EditionI | null | boolean> | EditionI | null | boolean
}

const emptyDefaultValues: EditionCreateInputI = {
  type: 'year',
  edition_name: '',
  tagline: '',
  description: '',
  registration_opens_at: '',
  registration_closes_at: '',
  starts_at: '',
  ends_at: '',
  timezone: typeof Intl.supportedValuesOf === 'function' ? Intl.supportedValuesOf('timeZone')[0] ?? 'UTC' : 'UTC',
  location_name: '',
  location_address: '',
  logo_url: '',
  banner_url: '',
  contact_email: '',
  contact_phone: '',
  organizer_name: '',
}

function toFormValues(edition: EditionI): EditionCreateInputI {
  return {
    type: edition.type,
    edition_name: edition.edition_name,
    tagline: edition.tagline ?? '',
    description: edition.description ?? '',
    registration_opens_at: edition.registration_opens_at ?? '',
    registration_closes_at: edition.registration_closes_at ?? '',
    starts_at: edition.starts_at,
    ends_at: edition.ends_at,
    timezone: edition.timezone,
    location_name: edition.location_name,
    location_address: edition.location_address,
    logo_url: edition.logo_url ?? '',
    banner_url: edition.banner_url ?? '',
    contact_email: edition.contact_email ?? '',
    contact_phone: edition.contact_phone ?? '',
    organizer_name: edition.organizer_name ?? '',
  }
}

export function ManageEditionModal({ open, onOpenChange, edition, onCreate, onUpdate }: ManageEditionModalProps) {
  const isEditing = Boolean(edition)
  const editValues = useMemo(() => (edition ? toFormValues(edition) : undefined), [edition])
  const steps = useMemo(() => createEditionFormSteps(), [])

  const controller = useMultiStepForm({
    schema: editionCreateSchema,
    steps,
    defaultValues: emptyDefaultValues,
    values: editValues,
    requireDirtyToSubmit: isEditing,
    onSubmit: async (values): Promise<boolean> => {
      const result = edition
        ? await onUpdate?.(edition.id, values)
        : await onCreate?.(values)

      return Boolean(result)
    },
    onSubmitSuccess: () => onOpenChange(false),
  })

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? 'Editar Edição' : 'Criar Nova Edição'}
      controller={controller}
      submitLabel={isEditing ? 'Salvar alterações' : 'Criar Edição'}
    />
  )
}
