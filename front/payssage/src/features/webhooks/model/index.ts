import z from "zod";

export const webhookEndpointCreateSchema = z.object({
  name: z.string().trim().min(1, "Name is required"),
  url: z.url({ error: "Invalid URL format" }),
});

export type WebhookEndpointCreateRequest = z.infer<
  typeof webhookEndpointCreateSchema
>;

export interface WebhookEndpoint {
  id: string;
  wallet_id: string;
  name: string;
  url: string;
  secret: string;
  created_at: string;
}

export interface WebhookEvent {
  id: string;
  wallet_id: string;
  intent_id: string;
  provider: string;
  external_id: string;
  event_type: string;
  payload: Record<string, unknown>;
  received_at: string;
}

export const webhookDeliveryStatuses = [
  "pending",
  "delivered",
  "failed",
] as const;
export type WebhookDeliveryStatus = (typeof webhookDeliveryStatuses)[number];

export interface WebhookDelivery {
  id: string;
  endpoint_id: string;
  event_id: string;
  status: WebhookDeliveryStatus;
  attempts: number;
  last_attempted_at: string | null;
  response_status: number | null;
  response_body: string | null;
  created_at: string;
}

// Compatibility aliases for the existing endpoint components.
export const webhookCreateSchema = webhookEndpointCreateSchema;
export type WebhookCreateI = WebhookEndpointCreateRequest;
export type WebhookI = WebhookEndpoint;
export type WebhookCreateResponseI = WebhookEndpoint;
