import { getFieldRulesStatus } from '@/shared/lib/forms/zod-utils'
import { Modal } from './modal'
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton'
import CrudForm from '@/shared/ui/form/CrudForm'
import type { CrudFormConfig, FieldConfig } from '@/shared/ui/form/types'
import type { FormValidateOrFn } from '@tanstack/react-form'
import type { ZodObject } from 'zod'

import type { FieldType, FieldOption } from '@/shared/ui/form/types'

export interface FieldDefinition {
  name: string
  label: string
  placeholder?: string
  required?: boolean
  type?: FieldType
  options?: FieldOption[]
}

interface FormModalProps<TFormData extends object> {
  isOpen: boolean
  onClose: () => void
  title: string
  description?: string
  schema: ZodObject
  fields: FieldDefinition[]
  onSubmit: (data: TFormData) => Promise<void> | void
  defaultValues: TFormData
  submitLabel?: string
  isLoading?: boolean
  formId: string
}

export function FormModal<TFormData extends object>({
  isOpen,
  onClose,
  title,
  description,
  schema,
  fields,
  onSubmit,
  defaultValues,
  submitLabel = 'Submit',
  isLoading,
  formId,
}: FormModalProps<TFormData>) {
  const validateFn: FormValidateOrFn<TFormData> = ({ value }) => {
    const result = schema.safeParse(value)
    if (!result.success) {
      const fieldErrors: Record<string, string> = {}
      for (const issue of result.error.issues) {
        const path = issue.path.join('.')
        fieldErrors[path] = issue.message
      }
      return { fields: fieldErrors }
    }
  }

  const fieldConfigs: FieldConfig[] = fields.map((f) => ({
    name: f.name,
    label: f.label,
    placeholder: f.placeholder ?? `Enter ${f.label.toLowerCase()}`,
    type: f.type ?? ('text' as const),
    required: f.required,
    getRulesStatus: (value) => getFieldRulesStatus(schema, f.name, value),
    options: f.options,
  }))

  const formOptions: CrudFormConfig<TFormData> = {
    defaultValues,
    validators: {
      onChange: validateFn,
      onSubmit: validateFn,
    },
    onSubmit: async ({ value }) => {
      await onSubmit(value)
    },
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      description={description}
      footer={
        <div className="flex justify-end pt-2 w-full">
          <ShadowButton
            type="submit"
            variant="accent-solid"
            formId={formId}
            disabled={isLoading}
            className="w-full rounded-sm font-bold transition-all h-10 justify-center"
            value={isLoading ? 'Submitting...' : submitLabel}
          />
        </div>
      }
    >
      <CrudForm<TFormData>
        key={formId}
        formId={formId}
        options={formOptions}
        fields={fieldConfigs}
      />
    </Modal>
  )
}