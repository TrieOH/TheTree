import z from 'zod'

const requiredDatetime = (label: string) =>
  z
    .string({ error: `${label} é obrigatório` })
    .min(1, { message: `${label} é obrigatório` })
    .pipe(z.iso.datetime({ message: `${label} inválido` }))

export const editionCreateSchema = z
  .object({
    name: z
      .string({ error: 'Nome é obrigatório' })
      .min(2, 'O nome deve ter pelo menos 2 caracteres')
      .max(256),
    slug: z
      .string({ error: 'Slug é obrigatório' })
      .min(2, 'O slug deve ter pelo menos 2 caracteres')
      .max(32),
    starts_at: requiredDatetime('Início'),
    ends_at: requiredDatetime('Término'),
  })
  .refine(
    ({ starts_at, ends_at }) =>
      new Date(ends_at).getTime() > new Date(starts_at).getTime(),
    {
      path: ['ends_at'],
      message: 'O término deve ser posterior ao início',
    },
  )

export type EditionCreateInputI = z.input<typeof editionCreateSchema>
export type EditionCreateOutputI = z.output<typeof editionCreateSchema>
export type EditionCreateSubmitI = EditionCreateOutputI

const nullableText = z
  .string()
  .transform((value) => value.trim() || null)
  .nullable()

const nullableDatetime = z
  .union([z.literal(''), z.iso.datetime({ message: 'Data inválida' })])
  .transform((value) => value || null)
  .nullable()

export const editionPatchSchema = editionCreateSchema.and(
  z.object({
    tagline: nullableText,
    description: nullableText,
    registration_opens_at: nullableDatetime,
    location_name: nullableText,
    location_description: nullableText,
    logo_url: nullableText,
    banner_url: nullableText,
    contact_email: z
      .union([z.literal(''), z.email('E-mail inválido')])
      .transform((value) => value || null)
      .nullable(),
  }),
)

export type EditionPatchInputI = z.input<typeof editionPatchSchema>
export type EditionPatchOutputI = z.output<typeof editionPatchSchema>

export type EditionStatus = 'draft' | 'future' | 'active' | 'past'

export interface EditionApiI {
  id: string
  event_id: string
  name: string
  slug: string
  tagline: string | null
  description: string | null
  is_draft: boolean
  registration_opens_at: string | null
  starts_at: string
  ends_at: string
  location_name: string | null
  location_description: string | null
  logo_url: string | null
  banner_url: string | null
  contact_email: string | null
  created_by: string
  created_at: string
  updated_at: string | null
  deleted_at: string | null
}

/**
 * UI model. `status` is always inferred from `is_draft` and the date range;
 * it is never expected from the API.
 */
export interface EditionI extends EditionApiI {
  status: EditionStatus
}

export function inferEditionStatus(
  edition: Pick<EditionApiI, 'is_draft' | 'starts_at' | 'ends_at'>,
  now = new Date(),
): EditionStatus {
  if (edition.is_draft) return 'draft'

  const currentTime = now.getTime()
  if (currentTime < new Date(edition.starts_at).getTime()) return 'future'
  if (currentTime > new Date(edition.ends_at).getTime()) return 'past'
  return 'active'
}

export function normalizeEdition(edition: EditionApiI): EditionI {
  return {
    ...edition,
    status: inferEditionStatus(edition),
  }
}

export function editionRangesOverlap(
  left: Pick<EditionI, 'starts_at' | 'ends_at'>,
  right: Pick<EditionCreateOutputI, 'starts_at' | 'ends_at'>,
) {
  return (
    new Date(left.starts_at).getTime() < new Date(right.ends_at).getTime() &&
    new Date(right.starts_at).getTime() < new Date(left.ends_at).getTime()
  )
}
