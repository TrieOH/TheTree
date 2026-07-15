import z from "zod";

const productTypeSchema = z
  .enum(
    ["merchandise", "ticket", "token", "bundle"],
    { error: "Tipo do Produto inválido" }
  ).default("merchandise");

export type ProductType = z.infer<typeof productTypeSchema>

const optionalNullableText = (schema: z.ZodType<string>) =>
  schema.optional().nullable().or(z.literal('')).transform((value) => (value === '' ? null : value))

const optionalNullableDatetime = () =>
  z.iso.datetime({ message: 'Data inválida' }).optional().nullable().or(z.literal('')).transform((value) => (value === '' ? null : value))

export const productCreateSchema = z.object({
  edition_scope_id: z.uuid(),
  name: z.string().min(3, 'Nome da edição precisa ter pelo menos 3 caracteres').max(256),
  description: optionalNullableText(z.string().max(8000)),
  type: productTypeSchema,
  ticket_id: z.preprocess(
    (val) => val === '' ? null : val,
    z.uuid({
      message: 'ID do ticket inválido',
    }).nullable().optional()
  ),
  price_cents: z.preprocess(
    (value) => (value === '' || value === null || value === undefined ? undefined : value),
    z.coerce.number({ error: 'O preço é obrigatório' }).int({
      message: 'O preço deve ser um número inteiro',
    }).nonnegative({
      message: 'O preço não pode ser negativo',
    }),
  ),
  available_from: optionalNullableDatetime(),
  available_until: optionalNullableDatetime(),
  thumbnail_url: optionalNullableText(z.url('Thumbnail deve ser uma URL válida')),
  gallery_urls: z.array(z.string()).nullish().transform(val => val ?? []),
  has_inventory: z.boolean().default(false),
  inventory_quantity: z.preprocess(
    (value) => (value === '' || value === null || value === undefined ? undefined : value),
    z.coerce.number({ error: 'Quantidade de estoque é obrigatória' }).int({ message: 'Quantidade de estoque precisa ser um número inteiro' }).nonnegative({ message: 'Quantidade de estoque não pode ser negativa' }),
  ).optional(),
}).superRefine((values, ctx) => {
  if (values.has_inventory && values.inventory_quantity === undefined) {
    ctx.addIssue({
      code: 'custom',
      path: ['inventory_quantity'],
      message: 'Quantidade de estoque é obrigatória',
    })
  }
})

export type ProductCreateInputI = z.input<typeof productCreateSchema>
export type ProductCreateOutputI = z.output<typeof productCreateSchema>

export const buyRequestItemSchema = z.object({
  product_id: z.uuid(),
  quantity: z.int().nonnegative().default(1)
})

export type BuyRequestItemI = z.infer<typeof buyRequestItemSchema>

export interface ProductI {
  id: string;
  scope_id: string;
  edition_id: string;
  name: string;
  description: string | null;
  type: ProductType;
  ticket_id: string | null;
  price_cents: number;
  status: "draft" | "available" | "sold_out" | "unavailable";
  available_from: string | null;
  available_until: string | null;
  has_inventory: boolean;
  inventory_quantity: number;
  inventory_remaining: number;
  created_by: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  thumbnail_url: string | null;
  gallery_urls: string[] | null;
}

export interface ReservedItemI {
  product_id: string;
  name: string;
  quantity: number;
  price_cents: number;
  product_type: ProductType;
  ticket_id?: string;
}

export interface UnavailableItemI {
  product_id: string;
  name: string;
  reason: string;
  requested: number;
  reserved: number;
}
