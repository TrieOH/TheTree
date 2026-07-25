import type { StepConfig } from '@/widgets/multi-step-form/model/types'
import type { CreateInitialProductInputI, ProductPatchInputI, VariantCreateInputI } from '.'

export function createProductFormSteps(): StepConfig<CreateInitialProductInputI>[] {
  return [
    {
      id: 'identidade',
      label: 'Identidade',
      fields: [
        {
          kind: 'text',
          name: 'vendor_code',
          label: 'Código do produto',
          placeholder: 'Ex: CAMISETA-UNIVENTS',
        },
        {
          kind: 'text',
          name: 'variant_vendor_code',
          label: 'Código da variação inicial',
          placeholder: 'Ex: CAMISETA-M',
        },
        {
          kind: 'text',
          name: 'name',
          label: 'Nome da variação',
          placeholder: 'Nome da primeira variação',
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
      ],
    },
    {
      id: 'comercial',
      label: 'Comercial',
      fields: [
        {
          kind: 'money',
          name: 'price',
          label: 'Preço',
          currency: "BRL",
          valueType: "number",
          maxCents: 99999999999,
        },
        {
          kind: 'toggle',
          name: 'requires_registration',
          label: 'Exige cadastro',
        },
      ],
    },
    {
      id: 'estoque',
      label: 'Estoque',
      fields: [
        {
          kind: 'text',
          name: 'stock',
          label: 'Quantidade em estoque',
          placeholder: 'Deixe vazio para ilimitado',
          inputType: 'number',
          optional: true,
        },
      ],
    },
  ]
}


export function createProductPatchFormSteps(): StepConfig<ProductPatchInputI>[] {
  return [
    {
      id: 'identidade',
      label: 'Identidade',
      fields: [
        {
          kind: 'text',
          name: 'vendor_code',
          label: 'Código do produto',
          placeholder: 'Ex: CAMISETA-UNIVENTS',
        },
        {
          kind: 'toggle',
          name: 'requires_registration',
          label: 'Exige cadastro',
        },
      ],
    },
  ]
}

export function createVariantFormSteps(): StepConfig<VariantCreateInputI>[] {
  return [
    {
      id: 'identidade',
      label: 'Identidade',
      fields: [
        {
          kind: 'text',
          name: 'vendor_code',
          label: 'Código do fornecedor',
          placeholder: 'Ex: CAMISETA-M-G',
        },
        {
          kind: 'text',
          name: 'name',
          label: 'Nome',
          placeholder: 'Nome da variação',
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
                placeholder="Descreva a variação"
                className="min-h-28 w-full rounded-xl border border-border/60 bg-background px-3 py-2.5 text-sm text-foreground outline-none transition-colors placeholder:text-muted-foreground/70 focus:border-primary focus:ring-2 focus:ring-primary/15"
              />
            </label>
          ),
        },
      ],
    },
    {
      id: 'comercial',
      label: 'Comercial',
      fields: [
        {
          kind: 'money',
          name: 'price',
          label: 'Preço',
          currency: "BRL",
          valueType: "number",
          maxCents: 99999999999,
        },
      ],
    },
    {
      id: 'estoque',
      label: 'Estoque',
      fields: [
        {
          kind: 'text',
          name: 'stock',
          label: 'Quantidade em estoque',
          placeholder: 'Deixe vazio para ilimitado',
          inputType: 'number',
          optional: true,
        },
      ],
    },
  ]
}