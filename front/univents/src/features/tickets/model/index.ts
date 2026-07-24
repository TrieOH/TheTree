import z from "zod"

export const ticketCreateSchema = z.object({
  name: z
    .string({ error: 'Nome é obrigatório' })
    .min(2, 'O nome deve ter pelo menos 2 caracteres.'),
  description: z.string().optional().nullable().transform((val) => (val === '' ? null : val)),
  access_level: z
    .coerce.number({ error: 'Nível de acesso é obrigatório' })
    .int({ message: 'Nível de acesso precisa ser um número inteiro' })
    .gte(0, { message: 'Nível de acesso não pode ser negativo' }),
  price_cents: z.preprocess(
    (value) => (value === '' || value === null || value === undefined ? undefined : value),
    z
      .coerce.number({ error: 'O preço é obrigatório' })
      .int({ message: 'O preço deve ser um número inteiro' })
      .nonnegative({ message: 'O preço não pode ser negativo' }),
  ),
  max_quantity: z.preprocess(
    (value) => (value === '' || value === null || value === undefined ? undefined : value),
    z
      .coerce.number({ error: 'Quantidade máxima é obrigatório' })
      .int({ message: 'Quantidade máxima precisa ser um número inteiro' })
      .gt(0, { message: 'Quantidade máxima precisa ser maior que zero' }),
  ).optional().nullable().or(z.literal('').transform(() => null)),
})

export type TicketCreateInputI = z.input<typeof ticketCreateSchema>
export type TicketCreateOutputI = z.output<typeof ticketCreateSchema>

export interface TicketI {
  id: string;
  edition_id: string;
  name: string;
  description: string | null;
  access_level: number;
  price_cents: number;
  max_quantity: number | null;
  created_by: string;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
}