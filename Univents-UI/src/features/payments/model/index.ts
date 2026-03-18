import z from "zod";

export const paymentConnectSchema = z.object({
  workspace_name: z.string().min(3),
  provider_redirect_url: z.url(),
  final_redirect_url: z.url(),
  provider: z.enum(["mercadopago"]), // FIXME: ADD Others
})

export type PaymentConnectI = z.infer<typeof paymentConnectSchema>

export const paymentDisconnectSchema = z.object({
  workspace_name: z.string().min(3),
  credential_id: z.string()
})

export type paymentDisconnectSchema = z.infer<typeof paymentConnectSchema>