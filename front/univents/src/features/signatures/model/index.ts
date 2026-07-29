import z from "zod";

export const signatureCreateSchema = z.object({
  title: z
    .string({ error: "Título é obrigatório" })
    .trim()
    .min(2, { message: "Título precisa ter pelo menos 2 caracteres" }),
  url: z.url("URL da assinatura inválida"),
});

export type SignatureCreateInputI = z.input<typeof signatureCreateSchema>;
export type SignatureCreateOutputI = z.output<typeof signatureCreateSchema>;

export interface SignatureI {
  id: string;
  edition_id: string;
  title: string;
  url: string;
}
