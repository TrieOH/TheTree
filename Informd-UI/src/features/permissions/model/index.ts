import z from "zod"

export const promoteToClientSchema = z.object({
  userId: z.string().min(1, "User ID is required"),
  requesterId: z.string().min(1, "Requester ID is required"),
})

export type PromoteToClientI = z.infer<typeof promoteToClientSchema>