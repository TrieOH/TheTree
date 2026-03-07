import { z } from "zod"

export const WebhookEndpointSchema = z.object({
    id: z.string().uuid(),
    workspace_id: z.string().uuid(),
    url: z.string().url(),
    created_at: z.string(),
})