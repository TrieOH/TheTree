import z from "zod"

export const imageURLUploadSchema = z.object({
  url: z.url(),
})

export type ImageURLUploadI = z.infer<typeof imageURLUploadSchema>