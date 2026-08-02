import type {
  CreateOAuthProviderRequest,
  OAuthProviderOutput,
  SupportedOAuthProviders,
  UpdateOAuthProviderRequest,
} from "@trieoh/identityx-api/schemas";
import z from "zod";

export const oauthProviderCreateSchema = z.object({
  provider: z.enum(["google", "github"]),
  client_id: z.string().min(1, "Client ID is required"),
  client_secret: z.string().min(1, "Client secret is required"),
}) satisfies z.ZodType<CreateOAuthProviderRequest>;

export const oauthProviderUpdateSchema = z.object({
  client_id: z.string().min(1, "Client ID is required"),
  client_secret: z.string().optional(),
}) satisfies z.ZodType<UpdateOAuthProviderRequest>;

export type OAuthProviderI = OAuthProviderOutput;
export type OAuthProvider = SupportedOAuthProviders;
export type OAuthProviderCreateI = z.input<typeof oauthProviderCreateSchema>;
export type OAuthProviderUpdateI = z.output<typeof oauthProviderUpdateSchema>;
