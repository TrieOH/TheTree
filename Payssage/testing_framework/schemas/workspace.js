import { z } from "zod"

export const WorkspaceSchema = z.object({
    id:         z.string().uuid(),
    name:       z.string(),
    created_at: z.string().datetime(),
})