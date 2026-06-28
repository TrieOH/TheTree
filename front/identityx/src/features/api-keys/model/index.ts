import type { ApiKey, CreateApiKeyResponse } from "@trieoh/identityx-models";
import z from "zod";

export const apiKeyCreateSchema = z.object({
  name: z.string().min(3, "Api key name must be at least 3 characters long"),
  create_for_service_account: z.enum(['true', 'false']).default('false'),
  expires_at: z.string().optional(),
});

export type ApiKeyCreateI = z.infer<typeof apiKeyCreateSchema>;

export type ApiKeyI = ApiKey;

export type CreateApiKeyResponseI = CreateApiKeyResponse;