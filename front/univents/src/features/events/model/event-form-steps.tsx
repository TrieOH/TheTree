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
