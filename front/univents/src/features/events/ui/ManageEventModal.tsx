import { useEffect, useRef } from 'react'
import { useMultiStepForm } from '@/widgets/multi-step-form/hooks/use-multi-step-form'
import { eventCreateSchema } from '../model'
import type { EventCreateInputI, EventCreateSubmitI, EventI } from '../model'
import { MultiStepFormModal } from '@/widgets/multi-step-form/ui/multi-step-form-modal'
import { createEventFormSteps } from '../model/event-form-steps'

export interface ManageEventModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreate: (
    values: EventCreateSubmitI,
  ) => Promise<EventI | null | boolean> | EventI | null | boolean
}

const emptyDefaultValues: EventCreateInputI = {
  full_name: '',
  slug: '',
  contact_email: '',
}

function toSlug(value: string): string {
  return value
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

function toAcronym(value: string): string {
  return value
    .trim()
    .split(/\s+/)
    .filter((word) => word.length > 2)
    .map((word) => word.charAt(0).toUpperCase())
    .join('')
}

export function ManageEventModal({
  open,
  onOpenChange,
  onCreate,
}: ManageEventModalProps) {
  const autoAcronymRef = useRef('')
  const autoSlugRef = useRef('')

  const controller = useMultiStepForm({
    schema: eventCreateSchema,
    steps: createEventFormSteps(),
    defaultValues: emptyDefaultValues,
    resetOnSuccessValues: emptyDefaultValues,
    onSubmit: async (values): Promise<boolean> => {
      // Media fields remain visible, but are intentionally excluded until
      // the API exposes routes for associating images with an event.
      const { logo_url: _logoUrl, banner_url: _bannerUrl, ...payload } = values
      return Boolean(await onCreate(payload))
    },
    onSubmitSuccess: () => onOpenChange(false),
  })

  const fullName = controller.form.watch('full_name')

  useEffect(() => {
    const trimmedName = fullName.trim()
    const nextAcronym = trimmedName ? toAcronym(trimmedName) : ''
    const nextSlug = trimmedName ? toSlug(trimmedName) : ''
    const currentAcronym = controller.form.getValues('acronym') ?? ''
    const currentSlug = controller.form.getValues('slug')

    if (currentAcronym === '' || currentAcronym === autoAcronymRef.current) {
      controller.form.setValue('acronym', nextAcronym, {
        shouldDirty: false,
        shouldTouch: false,
        shouldValidate: Boolean(trimmedName),
      })
    }

    if (currentSlug === '' || currentSlug === autoSlugRef.current) {
      controller.form.setValue('slug', nextSlug, {
        shouldDirty: false,
        shouldTouch: false,
        shouldValidate: Boolean(trimmedName),
      })
    }

    autoAcronymRef.current = nextAcronym
    autoSlugRef.current = nextSlug
  }, [controller.form, fullName])

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Criar Novo Evento"
      controller={controller}
      submitLabel="Criar Evento"
    />
  )
}
