import { useMemo } from 'react'
import { useMultiStepForm } from '@/widgets/multi-step-form/hooks/use-multi-step-form'
import { MultiStepFormModal } from '@/widgets/multi-step-form/ui/multi-step-form-modal'
import { productCreateSchema, type ProductCreateInputI, type ProductCreateOutputI, type ProductI } from '../model'
import { createProductFormSteps } from '../model/product-form-steps'

export interface ManageProductModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  editionId: string
  product?: ProductI
  onCreate?: (values: ProductCreateOutputI) => Promise<ProductI | null | boolean> | ProductI | null | boolean
  onUpdate?: (id: string, values: ProductCreateOutputI) => Promise<ProductI | null | boolean> | ProductI | null | boolean
}

const emptyDefaultValues: ProductCreateInputI = {
  edition_scope_id: '',
  name: '',
  description: '',
  type: 'merchandise',
  ticket_id: '',
  price_cents: 0,
  available_from: '',
  available_until: '',
  thumbnail_url: '',
  gallery_urls: [],
  has_inventory: false,
  inventory_quantity: 0,
}

function toFormValues(product: ProductI): ProductCreateInputI {
  return {
    edition_scope_id: product.scope_id,
    name: product.name,
    description: product.description ?? '',
    type: product.type,
    ticket_id: product.ticket_id ?? '',
    price_cents: product.price_cents,
    available_from: product.available_from ?? '',
    available_until: product.available_until ?? '',
    thumbnail_url: product.thumbnail_url ?? '',
    gallery_urls: product.gallery_urls ?? [],
    has_inventory: product.has_inventory,
    inventory_quantity: product.inventory_quantity,
  }
}

export function ManageProductModal({
  open,
  onOpenChange,
  product,
  editionId,
  onCreate,
  onUpdate,
}: ManageProductModalProps) {
  const isEditing = Boolean(product)
  const editValues = useMemo(() => (product ? toFormValues(product) : undefined), [product])
  const steps = useMemo(() => createProductFormSteps(), [])

  const controller = useMultiStepForm({
    schema: productCreateSchema,
    steps: steps,
    defaultValues: { ...emptyDefaultValues, edition_scope_id: editionId },
    values: editValues,
    requireDirtyToSubmit: isEditing,
    onSubmit: async (values): Promise<boolean> => {
      const result = product
        ? await onUpdate?.(product.id, values)
        : await onCreate?.(values)

      return Boolean(result)
    },
    onSubmitSuccess: () => onOpenChange(false),
  })

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? 'Editar produto' : 'Criar produto'}
      controller={controller}
      submitLabel={isEditing ? 'Salvar alterações' : 'Criar produto'}
    />
  )
}
