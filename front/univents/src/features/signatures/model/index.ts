import z from 'zod'

export const signatureCreateSchema = z.object({
  title: z.string().min(2),
  url: z.string().min(1),
  pos_x: z.number().int(),
  pos_y: z.number().int(),
})

export type SignatureCreateI = z.infer<typeof signatureCreateSchema>

export interface SignatureI {
  id: string
  edition_id: string
  title: string
  url: string
  pos_x: number
  pos_y: number
}
