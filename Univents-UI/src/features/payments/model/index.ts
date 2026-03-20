import z from "zod";

const paymentProviderSchema = z.enum(["mercadopago"]) // FIXME: ADD Others
export type PaymentProviderI = z.infer<typeof paymentProviderSchema>

export const paymentConnectSchema = z.object({
  workspace_name: z.string().min(3),
  provider_redirect_url: z.url(),
  final_redirect_url: z.url(),
  provider: paymentProviderSchema,
})

export type PaymentConnectI = z.infer<typeof paymentConnectSchema>

export const paymentDisconnectSchema = z.object({
  workspace_name: z.string().min(3),
  credential_id: z.string()
})

export type paymentDisconnectSchema = z.infer<typeof paymentConnectSchema>

export const submitPaymentPayload = z.object({
  card_token: z.string(),
  payment_method_id: z.string(),
  payment_method_type: z.string(),
  installments: z.int().nonnegative(),
  seller_credential_id: z.string(),
  payer_email: z.email(),
})

export type SubmitPaymentPayloadI = z.infer<typeof submitPaymentPayload>