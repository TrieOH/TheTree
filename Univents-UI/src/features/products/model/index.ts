import z from "zod";

const productTypeSchema = z
  .enum(
    ["merchandise", "ticket", "token", "bundle"],
    { error: "Invalid product type" }
  ).default("merchandise");

export type ProductType = z.infer<typeof productTypeSchema>

export const productCreateSchema = z.object({
  edition_scope_id: z.uuid(),
  name: z.string().min(3),
  description: z.string().optional().nullable(),
  type: productTypeSchema,
  ticket_id: z.uuid().optional().nullable(),
  price_cents: z.int().nonnegative(),
  available_from: z.iso.datetime().optional().nullable(),
  available_until: z.iso.datetime().optional().nullable(),
  has_inventory: z.boolean().default(false),
  inventory_quantity: z.int().nonnegative().default(0),
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

export const imageURLProductSchema = z.object({
  url: z.url(),
})

export type ImageURLProductI = z.infer<typeof imageURLProductSchema>

