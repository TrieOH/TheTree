import z from "zod";

const productTypeSchema = z
  .enum(
    ["merchandise", "ticket", "token", "bundle"],
    { error: "Tipo do Produto inválido" }
  ).default("merchandise");

export type ProductType = z.infer<typeof productTypeSchema>

export const productCreateSchema = z.object({
  edition_scope_id: z.uuid(),
  name: z.string({
    error: "Nome é obrigatório",
  }).min(3, {
    message: "O nome deve ter pelo menos 3 caracteres",
  }),
  description: z.string().nullable().optional(),
  type: productTypeSchema,
  ticket_id: z.preprocess(
    (val) => val === "" ? null : val,
    z.uuid({
      message: "ID do ticket inválido",
    }).nullable().optional()
  ),
  price_cents: z.int({
    message: "O preço deve ser um número inteiro",
  }).nonnegative({
    message: "O preço não pode ser negativo",
  }),
  available_from: z.preprocess(
    (val) => val === "" ? null : val,
    z.iso.datetime({
      message: "Data de início inválida",
    }).nullable().optional()
  ),
  available_until: z.preprocess(
    (val) => val === "" ? null : val,
    z.iso.datetime({
      message: "Data de término inválida",
    }).nullable().optional()
  ),
  thumbnail_url: z.string().optional().nullable().transform(val => val === "" ? null : val),
  gallery_urls: z.array(z.string()).nullish().transform(val => val ?? []),
  has_inventory: z.boolean().default(false),
  inventory_quantity: z.int({
    message: "Quantidade de estoque deve ser um número inteiro",
  }).nonnegative({
    message: "Quantidade de estoque não pode ser negativa",
  }).default(0),
})

export type ProductCreateI = z.infer<typeof productCreateSchema>

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
