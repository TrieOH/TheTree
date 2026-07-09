import { env } from "#/env";
import { createAppFetchers } from "@trieoh/api-client"

const { authFetcher, queryFetcher, publicFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_AUTH_API_URL,
  timeout: 10_000,
})

export { authFetcher, publicFetcher, queryFetcher as tanstackQueryFetcher }
