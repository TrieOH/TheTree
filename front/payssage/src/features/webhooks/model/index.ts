import z from "zod";

export const webhookCreateSchema = z.object({
  url: z.url({ error: "Invalid URL format" })
});

export type WebhookCreateI = z.infer<typeof webhookCreateSchema>;

export interface WebhookI {
  id: string;
  workspace_id: string;
  url: string;
  created_at: string;
}

export interface WebhookCreateResponseI extends WebhookI {
  secret: string;
}

