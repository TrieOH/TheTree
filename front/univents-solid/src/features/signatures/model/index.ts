import {
  Schema,
  createSchema,
  minLength,
  min,
  max,
  refine,
  type Infer,
} from "@/shared/lib/effect-schema";

export const signatureStep1Schema = createSchema(
  Schema.Struct({
    signatory_name: Schema.Trim.pipe(
      minLength(2, { message: "O nome do signatário deve ter pelo menos 2 caracteres." }),
    ),
    signatory_email: Schema.optional(
      Schema.Unknown.pipe(
        refine(
          (val) => {
            const str = String(val ?? "").trim();
            if (str.length === 0) return true;
            return str.includes("@") && str.includes(".");
          },
          { message: "Informe um e-mail válido." },
        ),
      ),
    ),
  }),
);

export const signatureCreateSchema = createSchema(
  Schema.Struct({
    signatory_name: Schema.Trim.pipe(
      minLength(2, { message: "O nome do signatário deve ter pelo menos 2 caracteres." }),
    ),
    signatory_title: Schema.optional(Schema.NullOr(Schema.String)),
    signatory_email: Schema.optional(Schema.NullOr(Schema.String)),
    signatory_user_id: Schema.optional(Schema.NullOr(Schema.String)),
    image_url: Schema.Trim.pipe(
      minLength(1, { message: "A URL da imagem de assinatura é obrigatória." }),
    ),
  }),
);

export type SignatureCreateInputI = Infer<typeof signatureCreateSchema>;
export type SignatureCreateOutputI = Infer<typeof signatureCreateSchema>;

export interface SignatureI {
  id: string;
  edition_id: string;
  created_by: string;
  signatory_name: string;
  signatory_title?: string | null;
  signatory_email?: string | null;
  signatory_user_id?: string | null;
  image_url: string;
  created_at: string;
  updated_at?: string | null;
  deleted_at?: string | null;
}

export type SignatureRequestStatus =
  | "pending"
  | "completed"
  | "expired"
  | "cancelled";

export interface SignatureRequestI {
  id: string;
  edition_id: string;
  created_by: string;
  signatory_name: string;
  signatory_title?: string | null;
  signatory_email?: string | null;
  signatory_user_id?: string | null;
  idempotency_key: string;
  status: SignatureRequestStatus;
  status_reason?: string | null;
  expires_at: string;
  signature_id?: string | null;
  created_at: string;
  updated_at?: string | null;
}

export interface SignatureRequestCreateI {
  signatory_name: string;
  signatory_title?: string;
  signatory_email: string;
  signatory_user_id?: string;
  expires_in_days?: number;
}

export const signatureRequestStep1Schema = createSchema(
  Schema.Struct({
    signatory_name: Schema.Trim.pipe(
      minLength(2, { message: "O nome do signatário deve ter pelo menos 2 caracteres." }),
    ),
    signatory_email: Schema.Unknown.pipe(
      refine(
        (val) => {
          const str = String(val ?? "").trim();
          return Boolean(str) && str.includes("@") && str.includes(".");
        },
        { message: "Informe um e-mail válido para envio do convite." },
      ),
    ),
    expires_in_days: Schema.Unknown.pipe(
      refine(
        (val) => {
          const days = Number(val);
          return !isNaN(days) && Number.isInteger(days) && days >= 1 && days <= 365;
        },
        { message: "A validade deve ser um número inteiro entre 1 e 365 dias." },
      ),
    ),
  }),
);

export const signatureRequestCreateSchema = createSchema(
  Schema.Struct({
    signatory_name: Schema.Trim.pipe(
      minLength(2, { message: "O nome do signatário deve ter pelo menos 2 caracteres." }),
    ),
    signatory_title: Schema.optional(Schema.String),
    signatory_email: Schema.Trim.pipe(
      minLength(3, { message: "Informe um e-mail válido para envio do convite." }),
    ),
    signatory_user_id: Schema.optional(Schema.String),
    expires_in_days: Schema.Number.pipe(
      min(1, { message: "A validade deve ser de pelo menos 1 dia." }),
      max(365, { message: "A validade máxima é de 365 dias." }),
    ),
  }),
);

export type SignatureRequestCreateInputI = Infer<
  typeof signatureRequestCreateSchema
>;
export type SignatureRequestCreateOutputI = Infer<
  typeof signatureRequestCreateSchema
>;

export * from "./token";
