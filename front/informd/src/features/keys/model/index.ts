import type {
  APIKeyResponse,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse
} from "@trieoh/informd-models";
import z from "zod";

export const apiKeyCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
}) satisfies z.ZodType<CreateAPIKeyRequest>;

export type ApiKeyCreateI = CreateAPIKeyRequest;

export type ApiKeyI = APIKeyResponse;
export type ApiKeyCreateResponseI = CreateAPIKeyResponse;
