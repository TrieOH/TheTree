import { z } from "zod"

export const APIKeySchema = z.object({
    id:         z.string().uuid(),
    name:       z.string(),
    prefix:     z.string(),
    created_at: z.string().datetime(),
    revoked_at: z.string().datetime().nullable(),
})

export const CreateAPIKeySchema = APIKeySchema.extend({
    key: z.string().startsWith("tp_"),
})