import type {
  ApiKey,
  CreateApiKeyRequest,
  CreateApiKeyResponse,
} from "@trieoh/identityx-models";
import z from "zod";

export const apiKeyCreateSchema = z.object({
  name: z.string().min(3, "Api key name must be at least 3 characters long"),
  capabilities: z.array(z.string()),
  subject_id: z
    .string()
    .min(3, "Subject ID must be at least 3 characters long")
    .optional(),
  env: z.string().min(3, "Environment must be at least 3 characters long"),
  expires_at: z.string().optional(),
}) satisfies z.ZodType<CreateApiKeyRequest>;

export type ApiKeyCreateI = z.infer<typeof apiKeyCreateSchema>;

export type ApiKeyI = ApiKey;

export type CreateApiKeyResponseI = CreateApiKeyResponse;
