import z from "zod";

export const signatureCreateSchema = z.object({
  signatory_name: z.string({ error: "Nome é obrigatório" }).trim().min(2),
  signatory_title: z.string().trim().optional(),
  signatory_email: z.email().optional(),
  signatory_user_id: z.uuid().optional(),
  image_url: z.url("URL da assinatura inválida"),
});

export type SignatureCreateInputI = z.input<typeof signatureCreateSchema>;
export type SignatureCreateOutputI = z.output<typeof signatureCreateSchema>;

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

export const signatureRequestCreateSchema = z.object({
  signatory_name: z.string().trim().min(2, "Informe o nome do signatário"),
  signatory_title: z.string().trim().optional(),
  signatory_email: z.email("Informe um e-mail válido"),
  signatory_user_id: z.string().uuid().optional().or(z.literal("")),
  expires_in_days: z.coerce
    .number()
    .int("Informe um número inteiro")
    .min(1, "O prazo mínimo é de 1 dia")
    .max(365, "O prazo máximo é de 365 dias"),
});

export type SignatureRequestCreateInputI = z.input<
  typeof signatureRequestCreateSchema
>;
export type SignatureRequestCreateOutputI = z.output<
  typeof signatureRequestCreateSchema
>;
