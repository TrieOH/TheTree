import {
  Schema,
  createSchema,
  min,
  max,
  minLength,
  pattern,
  nonNegative,
  type Infer,
} from "@/shared/lib/effect-schema";

const hexColorSchema = Schema.String.pipe(
  pattern(/^#([0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$/),
);

const elementBoundsFields = {
  id: Schema.String.pipe(minLength(1)),
  x: Schema.Number,
  y: Schema.Number,
  width: Schema.Number.pipe(min(1)),
  height: Schema.Number.pipe(min(1)),
};

const richTextRunSchema = Schema.Struct({
  text: Schema.String,
  bold: Schema.Boolean,
  italic: Schema.Boolean,
  underline: Schema.Boolean,
  color: hexColorSchema,
  fontSize: Schema.Number.pipe(min(6), max(400)),
  fontFamily: Schema.String.pipe(minLength(1)),
});

const richTextParagraphSchema = Schema.Struct({
  align: Schema.Literals(["left", "center", "right", "justify"]),
  lineHeight: Schema.Number.pipe(min(0.5), max(4)),
  runs: Schema.Array(richTextRunSchema),
});

const hashElementSchema = Schema.Struct({
  ...elementBoundsFields,
  type: Schema.Literal("hash"),
  hashLabel: Schema.String,
  hash: Schema.String,
  linkLabel: Schema.String,
  url: Schema.String,
  fontSize: Schema.Number.pipe(min(6), max(200)),
  color: hexColorSchema,
  align: Schema.Literals(["left", "center", "right"]),
});

const textElementSchema = Schema.Struct({
  ...elementBoundsFields,
  type: Schema.Literal("text"),
  paragraphs: Schema.Array(richTextParagraphSchema),
});

const imageElementSchema = Schema.Struct({
  ...elementBoundsFields,
  type: Schema.Literal("image"),
  src: Schema.String,
  fit: Schema.Literals(["cover", "contain", "fill"]),
  radius: Schema.Number.pipe(nonNegative(), max(400)),
  opacity: Schema.Number.pipe(nonNegative(), max(1)),
});

const signatureElementSchema = Schema.Struct({
  ...elementBoundsFields,
  type: Schema.Literal("signature"),
  signatureId: Schema.String.pipe(minLength(1)),
  src: Schema.String,
  name: Schema.String,
  fit: Schema.Literals(["cover", "contain", "fill"]),
  radius: Schema.Number.pipe(nonNegative(), max(400)),
  opacity: Schema.Number.pipe(nonNegative(), max(1)),
});

export const certificationTemplateElementSchema = createSchema(
  Schema.Union([
    hashElementSchema,
    textElementSchema,
    imageElementSchema,
    signatureElementSchema,
  ]),
);

export type CertificationTemplateElement =
  Infer<typeof certificationTemplateElementSchema>;

export const certificationTemplateKindSchema = createSchema(
  Schema.Literals(["edition_attendance", "program_attendance"]),
);

export const certificationTemplateDesignSchema = createSchema(
  Schema.Struct({
    canvas: Schema.optional(
      Schema.Struct({
        width: Schema.Number.pipe(min(320), max(6000)),
        height: Schema.Number.pipe(min(320), max(6000)),
      }),
    ),
    background: Schema.NullOr(Schema.String),
    elements: Schema.Array(certificationTemplateElementSchema),
  }),
);

export type CertificationTemplateDesignData =
  Infer<typeof certificationTemplateDesignSchema>;

export const certificationTemplateCreateSchema = createSchema(
  Schema.Struct({
    kind: certificationTemplateKindSchema,
    name: Schema.String.pipe(minLength(3)),
    description: Schema.optional(Schema.NullOr(Schema.String)),
    design_data: certificationTemplateDesignSchema,
  }),
);

export type CertificationTemplateCreateI =
  Infer<typeof certificationTemplateCreateSchema>;

export interface CertificationTemplateI {
  id: string;
  edition_id: string;
  kind: Infer<typeof certificationTemplateKindSchema>;
  name: string;
  description: string | null;
  design_data: Infer<typeof certificationTemplateDesignSchema>;
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
