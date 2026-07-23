import z from 'zod'

export const eventCreateSchema = z.object({
  full_name: z
    .string({ error: 'Nome é obrigatório' })
    .min(2, 'O nome deve ter pelo menos 2 caracteres.'),
  acronym: z
    .string()
    .optional()
    .nullable()
    .transform((val) => (val === '' ? null : val)),
  slug: z
    .string({ error: 'Slug é obrigatório' })
    .min(2, 'O slug deve ter pelo menos 2 caracteres.'),
  description: z
    .string()
    .optional()
    .nullable()
    .transform((val) => (val === '' ? null : val)),
  // Kept only as UI state until event media routes are available.
  logo_url: z
    .string()
    .optional()
    .nullable()
    .transform((val) => (val === '' ? null : val)),
  banner_url: z
    .string()
    .optional()
    .nullable()
    .transform((val) => (val === '' ? null : val)),
  contact_email: z
    .email('Informe um e-mail válido')
    .optional()
    .nullable()
    .or(z.literal(''))
    .transform((val) => (val === '' ? null : val)),
})

export type EventCreateInputI = z.input<typeof eventCreateSchema>
export type EventCreateOutputI = z.output<typeof eventCreateSchema>
export type EventCreateSubmitI = Omit<
  EventCreateOutputI,
  'logo_url' | 'banner_url'
>

export type EventStatusI = 'draft' | 'active' | 'discontinued'

export interface EventI {
  id: string
  owner_id: string
  full_name: string
  acronym: string | null
  slug: string
  description: string | null
  style: unknown | null
  payssage_seller_id: string | null
  payssage_wallet_id: string | null
  logo_url: string | null
  banner_url: string | null
  contact_email: string | null
  status: EventStatusI
  created_at: string
  updated_at: string | null
  deleted_at: string | null
  /** Legacy presentation fields; absent from the current events API. */
  tagline?: string | null
  is_series?: boolean
  editions_count: number
  gallery_urls?: string[] | null
  social_links?: {
    website?: string | null
    twitter?: string | null
    instagram?: string | null
    linkedin?: string | null
  } | null
}
