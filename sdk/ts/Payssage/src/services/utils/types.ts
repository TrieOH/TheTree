export interface ConnectRequestI {
  final_redirect_url: string;
  provider_redirect_url: string;
}

export interface IntentResponseI {
  id: string;
  workspace_id: string;
  amount: number;
  currency: string;
  status: "pending" | "succeeded" | "cancelled" | "failed";
  client_secret: string;
  provider: string;
  provider_payment_id: string | null;
  metadata: RawJsonValue;
  created_at: string;
  updated_at: string;
}

type RawJsonValue =
  | string
  | number
  | boolean
  | null
  | { [key: string]: RawJsonValue }
  | RawJsonValue[];