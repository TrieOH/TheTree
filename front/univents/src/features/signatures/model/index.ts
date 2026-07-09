import z from 'zod'

export const signatureCreateSchema = z.object({
  title: z.string().min(2),
  url: z.string().min(1),
})

export type SignatureCreateI = z.infer<typeof signatureCreateSchema>

export interface SignatureI {
  id: string
  edition_id: string
  title: string
  url: string
}
