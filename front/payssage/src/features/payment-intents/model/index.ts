import z from "zod";

export const intentStatuses = [
  "processing",
  "pending",
  "succeeded",
  "cancelled",
  "rejected",
  "failed",
  "refunded",
] as const;

export type IntentStatus = (typeof intentStatuses)[number];

export interface Intent {
  id: string;
  wallet_id: string;
  seller_id: string;
  collector_id: string;
  sandbox: boolean;
  amount_cents: number;
  currency: string;
  status: IntentStatus;
  provider: string;
  provider_data: Record<string, unknown>;
  metadata: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface CreateIntentRequest {
  seller_id: string;
  currency: string;
  amount_cents: number;
  checkout_provider_data: Record<string, unknown>;
}

export interface CreateIntentFormValues {
  seller_id: string;
  currency: string;
  amount_cents: number;
  checkout_provider_data: string;
}

export const createIntentSchema = z.object({
  seller_id: z.uuid("Seller ID must be a valid UUID"),
  currency: z
    .string()
    .trim()
    .length(3, "Currency must have exactly 3 characters")
    .transform((value) => value.toUpperCase()),
  amount_cents: z
    .number({ error: "Amount is required" })
    .int("Amount must be an integer")
    .positive("Amount must be greater than zero"),
  checkout_provider_data: z
    .string()
    .trim()
    .refine((value) => {
      try {
        const parsed: unknown = JSON.parse(value);
        return (
          typeof parsed === "object" &&
          parsed !== null &&
          !Array.isArray(parsed)
        );
      } catch {
        return false;
      }
    }, "Provider data must be a valid JSON object"),
}) satisfies z.ZodType<CreateIntentFormValues>;

export interface CreateTestModeIntentRequest {
  wallet_id: string;
  seller_id: string;
  collector_id?: string;
  amount_cents: number;
  currency: string;
  sandbox: boolean;
  provider: string;
  status: IntentStatus;
  provider_data: Record<string, unknown>;
  metadata?: Record<string, unknown>;
}

export interface CreateTestModeIntentFormValues {
  wallet_id: string;
  seller_id: string;
  collector_id: string;
  amount_cents: number;
  currency: string;
  sandbox: boolean;
  provider: string;
  status: IntentStatus;
  provider_data: string;
  metadata: string;
}

const jsonObjectString = (label: string) =>
  z
    .string()
    .trim()
    .refine((value) => {
      try {
        const parsed: unknown = JSON.parse(value);
        return (
          typeof parsed === "object" &&
          parsed !== null &&
          !Array.isArray(parsed)
        );
      } catch {
        return false;
      }
    }, `${label} must be a valid JSON object`);

export const createTestModeIntentSchema = z.object({
  wallet_id: z.uuid("Wallet must be selected"),
  seller_id: z.uuid("Seller must be selected"),
  collector_id: z.union([
    z.literal(""),
    z.uuid("Collector must be a valid UUID"),
  ]),
  amount_cents: z.number().int().positive("Amount must be greater than zero"),
  currency: z
    .string()
    .trim()
    .length(3, "Currency must have exactly 3 characters")
    .transform((value) => value.toUpperCase()),
  sandbox: z.boolean(),
  provider: z.string().trim().min(1, "Provider is required"),
  status: z.enum(intentStatuses),
  provider_data: jsonObjectString("Provider data"),
  metadata: jsonObjectString("Metadata"),
}) satisfies z.ZodType<CreateTestModeIntentFormValues>;
