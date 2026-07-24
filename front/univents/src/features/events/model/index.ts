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
  contact_email: z
    .email('Informe um e-mail válido')
    .optional()
    .nullable()
    .or(z.literal(''))
    .transform((val) => (val === '' ? null : val)),
})

export type EventCreateInputI = z.input<typeof eventCreateSchema>
export type EventCreateOutputI = z.output<typeof eventCreateSchema>

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
}
