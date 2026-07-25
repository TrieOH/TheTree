import { useMemo } from 'react'
import { useMultiStepForm } from '@/widgets/multi-step-form/hooks/use-multi-step-form'
import { MultiStepFormModal } from '@/widgets/multi-step-form/ui/multi-step-form-modal'
import { createInitialProductSchema, type CreateInitialProductInputI, type CreateInitialProductOutputI, type ProductI } from '../model'
import { createProductFormSteps } from '../model/product-form-steps'

export interface ManageProductModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  editionId: string
  onCreate: (values: CreateInitialProductOutputI) => Promise<ProductI | null | boolean>
}

const emptyDefaultValues: CreateInitialProductInputI = {
  requires_registration: false,
  vendor_code: '',
  variant_vendor_code: '',
  name: '',
  description: '',
  price: 0,
  stock: null,
}

export function ManageProductModal({
  open,
  onOpenChange,
  onCreate,
}: ManageProductModalProps) {
  const steps = useMemo(() => createProductFormSteps(), [])

  const controller = useMultiStepForm({
    schema: createInitialProductSchema,
    steps: steps,
    defaultValues: emptyDefaultValues,
    resetOnSuccessValues: emptyDefaultValues,
    onSubmit: async (values) => Boolean(await onCreate(values)),
    onSubmitSuccess: () => onOpenChange(false),
  })

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Criar Novo Produto"
      controller={controller}
      submitLabel="Criar Produto"
    />
  )
}