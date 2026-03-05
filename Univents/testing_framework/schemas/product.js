import { z } from "zod";

export const ProductSchema = z.object({
    id:                  z.string().uuid(),
    scope_id:            z.string().uuid(),
    edition_id:          z.string().uuid(),
    name:                z.string(),
    description:         z.string().nullable(),
    type:                z.enum(["merchandise", "ticket", "token", "bundle"]),
    ticket_id:           z.string().uuid().nullable(),
    price_cents:         z.number().int().nonnegative(),
    status:              z.enum(["draft", "available", "sold_out", "unavailable"]),
    available_from:      z.string().datetime().nullable(),
    available_until:     z.string().datetime().nullable(),
    has_inventory:       z.boolean(),
    inventory_quantity:  z.number().int().nonnegative(),
    inventory_remaining: z.number().int().nonnegative(),
    created_by:          z.string().uuid(),
    created_at:          z.string().datetime(),
    updated_at:          z.string().datetime(),
    deleted_at:          z.string().datetime().nullable(),
});