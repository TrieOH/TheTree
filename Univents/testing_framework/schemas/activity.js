import { z } from "zod";

export const ActivitySchema = z.object({
    id:                 z.string().uuid(),
    scope_id:           z.string().uuid(),
    edition_id:         z.string().uuid(),
    title:              z.string(),
    description:        z.string().nullable(),
    status:             z.string(),
    location:           z.string(),
    starts_at:          z.string().datetime(),
    ends_at:            z.string().datetime(),
    presenter_name:     z.string().nullable(),
    token_cost:         z.number().int(),
    has_capacity:       z.boolean(),
    capacity:           z.number().int(),
    remaining_capacity: z.number().int(),
    difficulty:         z.string().nullable(),
    created_by:         z.string().uuid(),
    created_at:         z.string().datetime(),
    updated_at:         z.string().datetime(),
    deleted_at:         z.string().datetime().nullable(),
});