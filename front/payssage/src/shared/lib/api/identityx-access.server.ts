import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts";
import { env } from "#/env";

export const identityXAccessClient = createIdentityXAccessClient({
  baseURL: env.VITE_AUTH_API_URL,
  apiKey: env.IDENTITYX_ACCESS_API_KEY,
});
