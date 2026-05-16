export interface PaymentIntentsI {
  id: string;
  workspace_id: string;
  amount: number;
  currency: string;
  status: "pending" | "succeeded" | "cancelled" | "failed";
  client_secret: string;
  provider: string;
  provider_payment_id: string | null;
  // metadata: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}