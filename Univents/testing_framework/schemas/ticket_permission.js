import { z } from "zod";

export const TicketPermissionSchema = z.object({
    id:              z.string().uuid(),
    ticket_id:       z.string().uuid(),
    permission_type: z.enum(["activity", "product", "checkpoint"]),
    activity_id:     z.string().uuid().nullable(),
    product_id:      z.string().uuid().nullable(),
    checkpoint_id:   z.string().uuid().nullable(),
    created_at:      z.string().datetime(),
});