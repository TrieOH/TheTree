import { env } from "@/env";
import { createAppFetchers } from "@trieoh/api-client"

const { authFetcher, authQueryFetcher, publicFetcher, publicQueryFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_AUTH_API_URL,
  timeout: 10_000,
})

export {
  authFetcher,
  authQueryFetcher,
  publicFetcher,
  publicQueryFetcher,
  authQueryFetcher as tanstackQueryFetcher,
}
