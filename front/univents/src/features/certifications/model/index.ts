import z from 'zod'

export const certificationTargetTypeSchema = z.enum(['edition', 'activity'], {
  error: 'Tipo de target inválido',
})

export type CertificationTargetType = z.infer<
  typeof certificationTargetTypeSchema
>

const hexColorSchema = z
  .string()
  .regex(/^#([0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$/, {
    message: 'Cor inválida, use um hexadecimal (ex: #1A2B3C)',
  })

const elementBoundsSchema = z.object({
  id: z.string().min(1),
  x: z.number(),
  y: z.number(),
  width: z.number().min(1),
  height: z.number().min(1),
})

const richTextRunSchema = z.object({
  text: z.string(),
  bold: z.boolean(),
  italic: z.boolean(),
  underline: z.boolean(),
  color: hexColorSchema,
  fontSize: z.number().min(6).max(400),
  fontFamily: z.string().min(1),
})

const richTextParagraphSchema = z.object({
  align: z.enum(['left', 'center', 'right', 'justify']),
  lineHeight: z.number().min(0.5).max(4).default(1.25),
  runs: z.array(richTextRunSchema),
})

export const certificationTemplateElementSchema = z.discriminatedUnion('type', [
  elementBoundsSchema.extend({
    type: z.literal('hash'),
    hashLabel: z.string(),
    hash: z.string(),
    linkLabel: z.string(),
    url: z.string(),
    fontSize: z.number().min(6).max(200),
    color: hexColorSchema,
    align: z.enum(['left', 'center', 'right']),
  }),
  elementBoundsSchema.extend({
    type: z.literal('text'),
    paragraphs: z.array(richTextParagraphSchema).min(1),
  }),
  elementBoundsSchema.extend({
    type: z.literal('image'),
    src: z.string(),
    fit: z.enum(['cover', 'contain', 'fill']),
    radius: z.number().min(0).max(400),
    opacity: z.number().min(0).max(1),
  }),
  elementBoundsSchema.extend({
    type: z.literal('signature'),
    signatureId: z.string().min(1),
    src: z.string(),
    name: z.string(),
    fit: z.enum(['cover', 'contain', 'fill']),
    radius: z.number().min(0).max(400),
    opacity: z.number().min(0).max(1),
  }),
])

export type CertificationTemplateElement = z.infer<
  typeof certificationTemplateElementSchema
>

export const certificationTemplateCreateSchema = z.object({
  title: z.string({ error: 'Título é obrigatório' }).min(3, {
    message: 'O título deve ter pelo menos 3 caracteres',
  }),
  url: z
    .string()
    .optional()
    .nullable()
    .transform((val) => (val === '' ? null : val)),
  data: z.object({
    background: z.string().nullable(),
    elements: z.array(certificationTemplateElementSchema),
  }),
})

export type CertificationTemplateCreateI = z.infer<
  typeof certificationTemplateCreateSchema
>

export interface CertificationTemplateI {
  id: string
  edition_id: string
  title: string
  data: z.infer<typeof certificationTemplateCreateSchema.shape.data>
  url: string | null // cert main background image url
  created_at: string
}

export interface CertificationI {
  id: string
  user_id: string
  target_id: string // edition or activity id
  target_type: CertificationTargetType
  certified_at: string
  hash: string
}

export interface VerifyCertificationResponseI {
  is_verified: boolean
  id: string
  user_id: string
  target_id: string
  target_type: CertificationTargetType
  certified_at: string
}

export interface SetCertificationTemplateRequestI {
  certification_template_id: string | null
}

export interface CertifyUserRequestI {
  target_id: string
  target_type: CertificationTargetType
}
