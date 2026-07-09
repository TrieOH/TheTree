import { env } from "#/env"
import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts"

export const identityXAccessClient = createIdentityXAccessClient({
  baseURL: env.VITE_AUTH_API_URL,
  apiKey: env.IDENTITYX_ACCESS_API_KEY,
})
