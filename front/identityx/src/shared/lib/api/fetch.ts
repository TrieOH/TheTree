import { createAppFetchers } from "@trieoh/api-client"
import { env } from "@/env"

const { authFetcher, queryFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_API_URL, // identityx uses same URL for both
  timeout: 10_000,
})

export { authFetcher, queryFetcher as tanstackQueryFetcher }
