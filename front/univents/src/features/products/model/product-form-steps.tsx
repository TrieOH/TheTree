import type { StepConfig } from '@/widgets/multi-step-form/model/types'
import type { ProductCreateInputI } from '.'

const productTypeOptions = [
  { label: 'Mercadoria', value: 'merchandise' },
  { label: 'Ingresso', value: 'ticket' },
  { label: 'Token', value: 'token' },
  { label: 'Pacote', value: 'bundle' },
]

export function createProductFormSteps(): StepConfig<ProductCreateInputI>[] {
  return [
    {
      id: 'identidade',
      label: 'Identidade',
      fields: [
        {
          kind: 'text',
          name: 'name',
          label: 'Nome',
          placeholder: 'Nome do produto',
        },
        {
          kind: 'custom',
          name: 'description',
          optional: true,
          render: ({ form }) => (
            <label className="block space-y-2">
              <span className="block text-sm font-medium text-foreground">Descrição</span>
              <textarea
                {...form.register('description')}
                rows={4}
                placeholder="Descreva o produto"
                className="min-h-28 w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm text-foreground outline-none transition-colors placeholder:text-muted-foreground/70 focus:border-primary focus:ring-2 focus:ring-primary/15"
              />
            </label>
          ),
        },
        {
          kind: 'combobox',
          name: 'type',
          label: 'Tipo',
          placeholder: 'Selecione o tipo',
          options: productTypeOptions,
        },
      ],
    },
    {
      id: 'comercial',
      label: 'Comercial',
      fields: [
        {
          kind: 'text',
          name: 'price_cents',
          label: 'Preço',
          placeholder: '0',
          inputType: 'number',
        },
        {
          kind: 'text',
          name: 'ticket_id',
          label: 'ID do ticket',
          placeholder: 'UUID do ticket',
          optional: true,
          visibleIf: { type: 'equals', field: 'type', value: 'ticket' },
        },
        {
          kind: 'datetime',
          name: 'available_from',
          label: 'Disponível de',
          optional: true,
        },
        {
          kind: 'datetime',
          name: 'available_until',
          label: 'Disponível até',
          optional: true,
        },
      ],
    },
    {
      id: 'estoque',
      label: 'Estoque',
      fields: [
        {
          kind: 'toggle',
          name: 'has_inventory',
          label: 'Controlar estoque',
        },
        {
          kind: 'text',
          name: 'inventory_quantity',
          label: 'Quantidade em estoque',
          placeholder: '0',
          inputType: 'number',
          optional: true,
          visibleIf: { type: 'equals', field: 'has_inventory', value: true },
        },
      ],
    },
    {
      id: 'midia',
      label: 'Mídia',
      fields: [
        {
          kind: 'image',
          name: 'thumbnail_url',
          label: 'Miniatura',
          hint: 'PNG ou JPEG, até 2MB',
          accept: 'image/png,image/jpeg',
          maxSizeMB: 2,
          optional: true,
        },
        {
          kind: 'gallery',
          name: 'gallery_urls',
          label: 'Galeria',
          hint: 'Até 6 imagens',
          accept: 'image/png,image/jpeg',
          maxSizeMB: 5,
          maxItems: 6,
          optional: true,
        },
      ],
    },
  ]
}
