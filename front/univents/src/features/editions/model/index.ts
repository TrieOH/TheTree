import z from "zod";

const editionTypeSchema = z
  .enum(["year", "season", "number", "ordinal", "custom"], { error: "Invalid edition type" }).default("year")
type EditionType = z.infer<typeof editionTypeSchema>

const optionalNullableText = (schema: z.ZodType<string>) =>
  schema.optional().nullable().or(z.literal('')).transform((value) => (value === '' ? null : value))

const optionalNullableDatetime = () =>
  z.iso.datetime({ message: 'Data inválida' }).optional().nullable().or(z.literal('')).transform((value) => (value === '' ? null : value))

const requiredDatetime = (label: string) =>
  z.string({ error: `${label} é obrigatória` })
    .min(1, { message: `${label} é obrigatória` })
    .pipe(z.iso.datetime({ message: `${label} inválida` }))

export const editionCreateSchema = z.object({
  type: editionTypeSchema,
  edition_name: z.string().min(3, 'Nome da edição precisa ter pelo menos 3 caracteres').max(256),
  tagline: optionalNullableText(z.string().max(512)),
  description: optionalNullableText(z.string().max(8000)),
  registration_opens_at: optionalNullableDatetime(),
  registration_closes_at: optionalNullableDatetime(),
  starts_at: requiredDatetime('Data de início'),
  ends_at: requiredDatetime('Data de término'),
  timezone: z.string().min(1, 'Fuso horário é obrigatório'),
  location_name: z.string().min(1, 'Nome do local é obrigatório'),
  location_address: z.string().min(1, 'Endereço é obrigatório'),
  logo_url: optionalNullableText(z.url('Logo deve ser uma URL válida')),
  banner_url: optionalNullableText(z.url('Banner deve ser uma URL válida')),
  contact_email: optionalNullableText(z.email('E-mail de contato deve ser válido')),
  contact_phone: optionalNullableText(z.string()),
  organizer_name: optionalNullableText(z.string()),
})

export type EditionCreateInputI = z.input<typeof editionCreateSchema>
export type EditionCreateOutputI = z.output<typeof editionCreateSchema>
export type EditionCreateSubmitI = Omit<EditionCreateOutputI, 'logo_url' | 'banner_url'>

export interface EditionI {
  id: string;
  event_id: string;
  goauth_scope_id: string;
  type: EditionType;
  edition_name: string;
  tagline: string | null;
  description: string | null;
  status: "draft" | "announced" | "open" | "ongoing" | "finished" |
  "completed" | "cancelled" | "postponed";
  monetary_type: "free" | "paid" | "mixed";
  registration_opens_at: string | null;
  registration_closes_at: string | null;
  starts_at: string;
  ends_at: string;
  timezone: string;
  location_name: string;
  location_address: string;
  logo_url: string | null;
  banner_url: string | null;
  contact_email: string | null;
  contact_phone: string | null;
  organizer_name: string | null;
  trie_payments_credential_id: string | null;
  trie_payments_provider: string | null;
  trie_payments_provider_public_key: string | null;
  certification_template_id?: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}
