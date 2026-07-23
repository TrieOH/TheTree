import type { StepConfig } from '@/widgets/multi-step-form/model/types'
import type { EventCreateInputI } from '.'

export function createEventFormSteps(): StepConfig<EventCreateInputI>[] {
  return [
    {
      id: 'identidade',
      label: 'Identidade',
      fields: [
        {
          kind: 'text',
          name: 'full_name',
          label: 'Nome do evento',
          placeholder: 'Ex: Tech Summit 2026',
        },
        {
          kind: 'text',
          name: 'slug',
          label: 'Slug',
          placeholder: 'tech-summit-2026',
        },
        {
          kind: 'text',
          name: 'acronym',
          label: 'Sigla',
          placeholder: 'TS26',
          optional: true,
        },
        {
          kind: 'text',
          name: 'description',
          label: 'Descrição',
          placeholder: 'Uma descrição curta sobre o evento',
          optional: true,
        },
      ],
    },
    {
      id: 'midia',
      label: 'Mídia',
      fields: [
        {
          kind: 'image',
          name: 'logo_url',
          label: 'Logo do evento',
          hint: 'PNG/JPEG (Máx. 2MB)',
          accept: 'image/png,image/jpeg',
          maxSizeMB: 2,
          optional: true,
          layout: 'half',
          disabled: true,
        },
        {
          kind: 'image',
          name: 'banner_url',
          label: 'Banner principal',
          hint: 'JPG/PNG (1920x1080)',
          accept: 'image/jpeg,image/png',
          maxSizeMB: 5,
          optional: true,
          layout: 'half',
          disabled: true,
        },
      ],
    },
    {
      id: 'conexao',
      label: 'Conexão',
      fields: [
        {
          kind: 'text',
          name: 'contact_email',
          label: 'E-mail de contato',
          placeholder: 'contato@evento.com',
          inputType: 'email',
        },
      ],
    },
  ]
}
