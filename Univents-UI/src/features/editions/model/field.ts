import type { FormFieldI } from "@/shared/model/field"
import type { EditionCreateI } from "."

const timezones = Intl.supportedValuesOf('timeZone')

export const getEditionFields = (): FormFieldI<EditionCreateI>[] => [
  {
    name: 'type' as const, label: 'Tipo', type: 'select' as const, options: [
      { value: 'year', label: 'Ano' },
      { value: 'season', label: 'Temporada' },
      { value: 'number', label: 'Número' },
      { value: 'ordinal', label: 'Ordinal' },
      { value: 'custom', label: 'Personalizado' },
    ], required: true
  },
  { name: 'edition_name' as const, label: 'Nome da edição', type: 'text' as const, placeholder: 'Ex: 2025', required: true },
  { name: 'starts_at' as const, label: 'Início', type: 'datetime' as const, required: true },
  { name: 'ends_at' as const, label: 'Fim', type: 'datetime' as const, required: true },
  { name: 'registration_opens_at' as const, label: 'Abertura de inscrições', type: 'datetime' as const },
  { name: 'registration_closes_at' as const, label: 'Fechamento de inscrições', type: 'datetime' as const },
  { name: 'location_name' as const, label: 'Nome do local', type: 'text' as const, placeholder: 'Ex: Centro de Convenções', required: true },
  { name: 'location_address' as const, label: 'Endereço', type: 'text' as const, placeholder: 'Rua, número, cidade', required: true },
  { name: 'timezone' as const, label: 'Fuso horário', type: 'select' as const, options: timezones.map(tz => ({ value: tz, label: tz })), required: true },
  { name: 'tagline' as const, label: 'Tagline', type: 'text' as const },
  { name: 'organizer_name' as const, label: 'Organizador', type: 'text' as const },
]