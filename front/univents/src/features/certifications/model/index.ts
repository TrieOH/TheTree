import z from "zod";

export const certificationTargetTypeSchema = z.enum(["edition", "activity"], {
  error: "Tipo de target inválido",
});

export type CertificationTargetType = z.infer<typeof certificationTargetTypeSchema>;

export const certificationTemplateCreateSchema = z.object({
  title: z.string({ error: "Título é obrigatório" }).min(3, {
    message: "O título deve ter pelo menos 3 caracteres",
  }),
  url: z.string().optional().nullable().transform((val) => val === "" ? null : val),
  data: z.object({
    background: z.string().nullable(),
    elements: z.array(z.discriminatedUnion('type', [
      z.object({
        type: z.literal('text'),
        xPct: z.number(),
        yPct: z.number(),
        widthPct: z.number(),
        heightPct: z.number(),
        content: z.string(),
      }),
      z.object({
        type: z.literal('signature'),
        xPct: z.number(),
        yPct: z.number(),
        widthPct: z.number(),
        heightPct: z.number(),
        title: z.string().nullable(),
        signatureId: z.string().nullable(),
      }),
      z.object({
        type: z.literal('image'),
        xPct: z.number(),
        yPct: z.number(),
        widthPct: z.number(),
        heightPct: z.number(),
        src: z.string().nullable(),
        fileName: z.string().nullable(),
      }),
    ])),
  }),
});

export type CertificationTemplateCreateI = z.infer<typeof certificationTemplateCreateSchema>;

export interface CertificationTemplateI {
  id: string;
  edition_id: string;
  title: string;
  data: z.infer<typeof certificationTemplateCreateSchema.shape.data>;
  url: string | null; // cert main background image url
  created_at: string;
}

export interface CertificationI {
  id: string;
  user_id: string;
  target_id: string; // edition or activity id
  target_type: CertificationTargetType;
  certified_at: string;
  hash: string;
}

export interface VerifyCertificationResponseI {
  is_verified: boolean;
  id: string;
  user_id: string;
  target_id: string;
  target_type: CertificationTargetType;
  certified_at: string;
}

export interface SetCertificationTemplateRequestI {
  certification_template_id: string | null;
}

export interface CertifyUserRequestI {
  target_id: string;
  target_type: CertificationTargetType;
}
