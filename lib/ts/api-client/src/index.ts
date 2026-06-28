/**
 * @trieoh/api-client
 *
 * Unified fetcher factory for frontend apps.
 *
 * Usage:
 *   import { createAppFetchers } from "@trieoh/api-client"
 *
 *   const { authFetcher, queryFetcher, publicFetcher } = createAppFetchers({
 *     apiURL: env.VITE_API_URL,
 *     authAPIURL: env.VITE_AUTH_API_URL,
 *   })
 */
import {
  createFetcher,
  createQueryFetcher,
  type ApiError,
} from "@trieoh/identityx-sdk-ts"
import { createDefaultFetchClient } from "@trieoh/envoy-fetch-ts"

export type { ApiError }

export interface AppFetcherConfig {
  /** Base URL for the main API. */
  apiURL: string
  /** Base URL for the authentication API (may differ from apiURL). */
  authAPIURL?: string
  /** Timeout in milliseconds (default: 10_000). */
  timeout?: number
}

export interface AppFetchers {
  /** Authenticated fetcher (auto-attaches Bearer token, handles refresh). */
  authFetcher: ReturnType<typeof createFetcher>
  /** Fetcher that returns raw data (for TanStack Query). */
  queryFetcher: ReturnType<typeof createQueryFetcher>
  /** Public/unauthenticated fetcher (no auth headers). */
  publicFetcher: ReturnType<typeof createDefaultFetchClient>
}

/**
 * Create all three fetcher variants every app needs.
 */
export function createAppFetchers(config: AppFetcherConfig): AppFetchers {
  const {
    apiURL,
    authAPIURL = apiURL,
    timeout = 10_000,
  } = config

  const clientConfig = { timeout }

  const authFetcher = createFetcher({
    baseURL: apiURL,
    authBaseURL: authAPIURL,
    clientConfig,
  })

  const queryFetcher = createQueryFetcher({
    baseURL: apiURL,
    authBaseURL: authAPIURL,
    clientConfig,
  })

  const publicFetcher = createDefaultFetchClient({
    baseURL: apiURL,
    credentials: "omit",
    timeout,
  })

  return { authFetcher, queryFetcher, publicFetcher }
}
