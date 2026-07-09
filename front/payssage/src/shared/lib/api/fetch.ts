import { env } from "#/env";
import { createAppFetchers } from "@trieoh/api-client"
import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts"

const { authFetcher, queryFetcher, publicFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_AUTH_API_URL,
  timeout: 10_000,
})

const identityXAccessClient = createIdentityXAccessClient({
  baseURL: env.VITE_AUTH_API_URL,
  apiKey: env.VITE_IDENTITYX_ACCESS_API_KEY,
})

export { authFetcher, publicFetcher, queryFetcher as tanstackQueryFetcher, identityXAccessClient }
