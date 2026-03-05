import { z } from "zod";

export const TicketSchema = z.object({
    id:          z.string().uuid(),
    scope_id:    z.string().uuid(),
    edition_id:  z.string().uuid(),
    name:        z.string(),
    description: z.string().nullable(),
    created_by:  z.string().uuid(),
    created_at:  z.string().datetime(),
    updated_at:  z.string().datetime(),
    deleted_at:  z.string().datetime().nullable(),
});