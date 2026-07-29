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
