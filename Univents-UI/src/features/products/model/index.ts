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
}