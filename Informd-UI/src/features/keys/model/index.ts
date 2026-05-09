import z from "zod";

export const apiKeyCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
});

export type ApiKeyCreateI = z.infer<typeof apiKeyCreateSchema>;

export interface ApiKeyI {
  id: string;
  name: string;
  prefix: string;
  created_at: string;
  revoked_at: string;
}

export interface ApiKeyCreateResponseI extends ApiKeyI {
  key: string;
}