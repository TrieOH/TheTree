import z from "zod";

const hexColorSchema = z
  .string()
  .regex(/^#([0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$/, {
    message: "Cor inválida, use um hexadecimal (ex: #1A2B3C)",
  });

const elementBoundsSchema = z.object({
  id: z.string().min(1),
  x: z.number(),
  y: z.number(),
  width: z.number().min(1),
  height: z.number().min(1),
});

const richTextRunSchema = z.object({
  text: z.string(),
  bold: z.boolean(),
  italic: z.boolean(),
  underline: z.boolean(),
  color: hexColorSchema,
  fontSize: z.number().min(6).max(400),
  fontFamily: z.string().min(1),
});

const richTextParagraphSchema = z.object({
  align: z.enum(["left", "center", "right", "justify"]),
  lineHeight: z.number().min(0.5).max(4).default(1.25),
  runs: z.array(richTextRunSchema),
});

export const certificationTemplateElementSchema = z.discriminatedUnion("type", [
  elementBoundsSchema.extend({
    type: z.literal("hash"),
    hashLabel: z.string(),
    hash: z.string(),
    linkLabel: z.string(),
    url: z.string(),
    fontSize: z.number().min(6).max(200),
    color: hexColorSchema,
    align: z.enum(["left", "center", "right"]),
  }),
  elementBoundsSchema.extend({
    type: z.literal("text"),
    paragraphs: z.array(richTextParagraphSchema).min(1),
  }),
  elementBoundsSchema.extend({
    type: z.literal("image"),
    src: z.string(),
    fit: z.enum(["cover", "contain", "fill"]),
    radius: z.number().min(0).max(400),
    opacity: z.number().min(0).max(1),
  }),
  elementBoundsSchema.extend({
    type: z.literal("signature"),
    signatureId: z.string().min(1),
    src: z.string(),
    name: z.string(),
    fit: z.enum(["cover", "contain", "fill"]),
    radius: z.number().min(0).max(400),
    opacity: z.number().min(0).max(1),
  }),
]);

export type CertificationTemplateElement = z.infer<
  typeof certificationTemplateElementSchema
>;

export const certificationTemplateKindSchema = z.enum([
  "edition_attendance",
  "program_attendance",
]);

export const certificationTemplateDesignSchema = z.object({
  canvas: z
    .object({
      width: z.number().min(320).max(6000),
      height: z.number().min(320).max(6000),
    })
    .optional(),
  background: z.string().nullable(),
  elements: z.array(certificationTemplateElementSchema),
});

export const certificationTemplateCreateSchema = z.object({
  kind: certificationTemplateKindSchema,
  name: z.string({ error: "Nome é obrigatório" }).min(3, {
    message: "O título deve ter pelo menos 3 caracteres",
  }),
  description: z.string().nullable().optional(),
  design_data: certificationTemplateDesignSchema,
});

export type CertificationTemplateCreateI = z.infer<
  typeof certificationTemplateCreateSchema
>;

export interface CertificationTemplateI {
  id: string;
  edition_id: string;
  kind: z.infer<typeof certificationTemplateKindSchema>;
  name: string;
  description: string | null;
  design_data: z.infer<typeof certificationTemplateDesignSchema>;
  created_at: string;
}

export interface CertificationTemplateProgramI {
  id: string;
  template_id: string;
  program_id: string;
  created_at: string;
}

export interface CertificationI {
  id: string;
  edition_id: string;
  template_id: string | null;
  registration_id: string;
  user_id: string;
  program_id: string | null;
  verification_hash: string;
  valid: boolean;
  invalid_reason: string | null;
  email_sent: boolean;
  issued_at: string;
  created_at: string;
  updated_at: string | null;
}

export interface CertificationEmissionErrorI {
  id: string;
  edition_id: string;
  user_id: string;
  template_id: string | null;
  program_id: string | null;
  error_message: string;
  created_at: string;
}

export interface VerifyCertificationResponseI {
  valid: boolean;
  template_id: string | null;
  cert: CertificationI | null;
}
