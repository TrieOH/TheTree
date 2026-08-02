/**
 * @trieoh/api-client
 *
 * Unified fetcher factory for frontend apps.
 *
 * Usage:
 *   import { createAppFetchers } from "@trieoh/api-client"
 *
 *   const { authFetcher, authQueryFetcher, publicFetcher, publicQueryFetcher } = createAppFetchers({
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
export {
  configureApiClient,
  createOrvalTransport,
  type ApiClientConfig,
  type ApiTransport,
} from "./orval-mutator"

/** Extracts the payload from Orval's full HTTP response wrapper. */
export const orvalData = <T>(response: { data: unknown }): T =>
  response.data as T

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
  /** Authenticated fetcher for TanStack Query. */
  authQueryFetcher: ReturnType<typeof createQueryFetcher>
  /** Public/unauthenticated fetcher (no auth headers). */
  publicFetcher: ReturnType<typeof createDefaultFetchClient>
  /** Public fetcher for TanStack Query. */
  publicQueryFetcher: <T>(path: string) => Promise<T>
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

  const authQueryFetcher = createQueryFetcher({
    baseURL: apiURL,
    authBaseURL: authAPIURL,
    clientConfig,
  })

  const publicFetcher = createDefaultFetchClient({
    baseURL: apiURL,
    credentials: "omit",
    timeout,
  })

  const publicQueryFetcher = async <T>(path: string): Promise<T> => {
    const result = await publicFetcher.get<T>(path)
    if (!result.success) throw result
    return result.data
  }

  return { authFetcher, authQueryFetcher, publicFetcher, publicQueryFetcher }
}
