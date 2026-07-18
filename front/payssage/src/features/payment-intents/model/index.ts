import z from "zod"

export const intentStatuses = ["processing", "succeeded", "cancelled", "failed"] as const

export type IntentStatus = (typeof intentStatuses)[number]

export interface Intent {
  id: string
  wallet_id: string
  seller_id: string
  collector_id: string
  sandbox: boolean
  amount_cents: number
  currency: string
  status: IntentStatus
  provider: string
  provider_data: Record<string, unknown>
  metadata: Record<string, unknown>
  created_at: string
  updated_at: string
}

export const intentSchema = z.object({
  id: z.string(),
  wallet_id: z.string(),
  seller_id: z.string(),
  collector_id: z.string(),
  sandbox: z.boolean(),
  amount_cents: z.number().int().positive("Amount must be greater than zero"),
  currency: z.string().length(3),
  status: z.enum(intentStatuses),
  provider: z.string(),
  provider_data: z.record(z.string(), z.unknown()),
  metadata: z.record(z.string(), z.unknown()),
  created_at: z.string(),
  updated_at: z.string(),
}) satisfies z.ZodType<Intent>

export type PaymentIntentsI = Intent
