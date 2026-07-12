import type { ComboboxOption, StepConfig } from '@/widgets/multi-step-form/model/types'
import type { EditionCreateInputI } from '.'

const timezones = typeof Intl.supportedValuesOf === 'function'
  ? Intl.supportedValuesOf('timeZone')
  : ['UTC']

const editionTypeOptions: ComboboxOption[] = [
  { value: 'year', label: 'Ano' },
  { value: 'season', label: 'Temporada' },
  { value: 'number', label: 'Número' },
  { value: 'ordinal', label: 'Ordinal' },
  { value: 'custom', label: 'Personalizado' },
]

const timezoneOptions: ComboboxOption[] = timezones.map((tz) => ({ value: tz, label: tz }))

export function createEditionFormSteps(): StepConfig<EditionCreateInputI>[] {
  return [
    {
      id: 'identidade',
      label: 'Identidade',
      fields: [
        {
          kind: 'combobox',
          name: 'type',
          label: 'Tipo',
          placeholder: 'Selecione o tipo',
          options: editionTypeOptions,
        },
        {
          kind: 'text',
          name: 'edition_name',
          label: 'Nome da edição',
          placeholder: 'Ex: 2025',
        },
        {
          kind: 'text',
          name: 'tagline',
          label: 'Tagline',
          placeholder: 'Uma descrição curta da edição',
          optional: true,
        },
        {
          kind: 'text',
          name: 'organizer_name',
          label: 'Organizador',
          placeholder: 'Nome do organizador',
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
          label: 'Logo da edição',
          hint: 'PNG ou JPEG, até 2MB',
          accept: 'image/png,image/jpeg',
          maxSizeMB: 2,
          optional: true,
          layout: 'half',
        },
        {
          kind: 'image',
          name: 'banner_url',
          label: 'Banner da edição',
          hint: 'PNG ou JPEG, até 5MB',
          accept: 'image/png,image/jpeg',
          maxSizeMB: 5,
          optional: true,
          layout: 'half',
        },
      ],
    },
    {
      id: 'periodo',
      label: 'Período',
      fields: [
        {
          kind: 'datetime',
          name: 'starts_at',
          label: 'Início',
          placeholder: 'Selecione a data de início',
        },
        {
          kind: 'datetime',
          name: 'ends_at',
          label: 'Fim',
          placeholder: 'Selecione a data de término',
        },
        {
          kind: 'datetime',
          name: 'registration_opens_at',
          label: 'Abertura de inscrições',
          placeholder: 'Opcional',
          optional: true,
        },
        {
          kind: 'datetime',
          name: 'registration_closes_at',
          label: 'Fechamento de inscrições',
          placeholder: 'Opcional',
          optional: true,
        },
      ],
    },
    {
      id: 'local',
      label: 'Local',
      fields: [
        {
          kind: 'text',
          name: 'location_name',
          label: 'Nome do local',
          placeholder: 'Ex: Centro de Convenções',
        },
        {
          kind: 'text',
          name: 'location_address',
          label: 'Endereço',
          placeholder: 'Rua, número, cidade',
        },
        {
          kind: 'combobox',
          name: 'timezone',
          label: 'Fuso horário',
          placeholder: 'Selecione o fuso',
          options: timezoneOptions,
        },
      ],
    },
    {
      id: 'contato',
      label: 'Contato',
      fields: [
        {
          kind: 'text',
          name: 'contact_email',
          label: 'E-mail de contato',
          placeholder: 'contato@evento.com',
          inputType: 'email',
          optional: true,
        },
        {
          kind: 'text',
          name: 'contact_phone',
          label: 'Telefone de contato',
          placeholder: '(00) 00000-0000',
          inputType: 'tel',
          optional: true,
        },
      ],
    },
  ]
}
