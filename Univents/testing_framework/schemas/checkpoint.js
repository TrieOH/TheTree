import { z } from "zod";

export const CheckpointSchema = z.object({
    id:          z.string().uuid(),
    scope_id:    z.string().uuid(),
    edition_id:  z.string().uuid(),
    name:        z.string(),
    type:        z.enum(["entry", "zone", "amenity", "session", "exit"]),
    access_mode: z.enum(["open", "ticket", "staff_only"]),
    starts_at:   z.string().datetime(),
    ends_at:     z.string().datetime(),
    created_by:  z.string().uuid(),
    created_at:  z.string().datetime(),
    updated_at:  z.string().datetime(),
    deleted_at:  z.string().datetime().nullable(),
});