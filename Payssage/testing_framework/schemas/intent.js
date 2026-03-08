import { z } from "zod"

export const IntentSchema = z.object({
    id:                  z.string().uuid(),
    workspace_id:        z.string().uuid(),
    amount:              z.number().int().positive(),
    currency:            z.string().length(3),
    status:              z.enum(["pending", "succeeded", "cancelled", "failed"]),
    client_secret:       z.string(),
    provider:            z.string(),
    provider_payment_id: z.string().nullable().optional(),
    metadata:            z.record(z.string(), z.any()).nullable(),
    created_at:          z.string().datetime(),
    updated_at:          z.string().datetime(),
})