import z from "zod";

export const createInitialProductSchema = z.object({
  requires_registration: z.boolean().default(false),
  vendor_code: z
    .string()
    .min(2, "Código precisa ter pelo menos 2 caracteres")
    .max(255),
  variant_vendor_code: z
    .string()
    .min(2, "Código da variação precisa ter pelo menos 2 caracteres")
    .max(255),
  name: z.string().min(2, "Nome precisa ter pelo menos 2 caracteres"),
  description: z
    .string()
    .max(8000)
    .optional()
    .nullable()
    .or(z.literal(""))
    .transform((v) => (v === "" ? null : v)),
  price: z.preprocess(
    (value) =>
      value === "" || value === null || value === undefined ? undefined : value,
    z.coerce
      .number({ error: "O preço é obrigatório" })
      .int({ message: "O preço deve ser um número inteiro" })
      .nonnegative({ message: "O preço não pode ser negativo" }),
  ),
  stock: z
    .preprocess(
      (value) =>
        value === "" || value === null || value === undefined
          ? undefined
          : value,
      z.coerce
        .number({ error: "Quantidade máxima é obrigatório" })
        .int({ message: "Quantidade máxima precisa ser um número inteiro" })
        .gt(0, { message: "Quantidade máxima precisa ser maior que zero" }),
    )
    .optional()
    .nullable()
    .or(z.literal("").transform(() => null)),
});

export type CreateInitialProductInputI = z.input<
  typeof createInitialProductSchema
>;
export type CreateInitialProductOutputI = z.output<
  typeof createInitialProductSchema
>;

export const productPatchSchema = z.object({
  requires_registration: z.boolean().default(false),
  vendor_code: z
    .string()
    .min(2, "Código precisa ter pelo menos 2 caracteres")
    .max(255),
});

export type ProductPatchInputI = z.input<typeof productPatchSchema>;
export type ProductPatchOutputI = z.output<typeof productPatchSchema>;

export interface ProductI {
  id: string;
  edition_id: string;
  vendor_code: string;
  requires_registration: boolean;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
}

export const variantCreateSchema = z.object({
  vendor_code: z
    .string()
    .min(2, "Código precisa ter pelo menos 2 caracteres")
    .max(255),
  name: z.string().min(2, "Nome precisa ter pelo menos 2 caracteres"),
  description: z
    .string()
    .max(8000)
    .optional()
    .nullable()
    .or(z.literal(""))
    .transform((v) => (v === "" ? null : v)),
  price: z.preprocess(
    (value) =>
      value === "" || value === null || value === undefined ? undefined : value,
    z.coerce
      .number({ error: "O preço é obrigatório" })
      .int({ message: "O preço deve ser um número inteiro" })
      .nonnegative({ message: "O preço não pode ser negativo" }),
  ),
  stock: z
    .preprocess(
      (value) =>
        value === "" || value === null || value === undefined
          ? undefined
          : value,
      z.coerce
        .number({ error: "Quantidade máxima é obrigatório" })
        .int({ message: "Quantidade máxima precisa ser um número inteiro" })
        .gt(0, { message: "Quantidade máxima precisa ser maior que zero" }),
    )
    .optional()
    .nullable()
    .or(z.literal("").transform(() => null)),
  gallery_urls: z.array(z.string()).default([]),
});

export type VariantCreateInputI = z.input<typeof variantCreateSchema>;
export type VariantCreateOutputI = z.output<typeof variantCreateSchema>;

export interface VariantI {
  id: string;
  edition_id: string;
  product_id: string;
  vendor_code: string;
  name: string;
  description: string | null;
  price: number;
  stock: number | null;
  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
  gallery_urls: string[];
}

export const buyRequestItemSchema = z.object({
  product_id: z.uuid(),
  quantity: z.int().nonnegative().default(1),
});

export type BuyRequestItemI = z.infer<typeof buyRequestItemSchema>;

export interface ReservedItemI {
  product_id: string;
  name: string;
  quantity: number;
  price_cents: number;
  product_type: string;
  ticket_id?: string;
}

export interface UnavailableItemI {
  product_id: string;
  name: string;
  reason: string;
  requested: number;
  reserved: number;
}
