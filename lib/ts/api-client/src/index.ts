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
import { envelopeFetch, itemsOf } from "./envelope"

export type { ApiError }
export { envelopeFetch, isEnvelope, itemsOf, toEnvelope, toEnvelopeResponse, type WireEnvelope } from "./envelope"
export {
  configureApiClient,
  createOrvalTransport,
  customInstance,
  type ApiClientConfig,
  type ApiTransport,
} from "./orval-mutator"

/** Extracts the payload from Orval's full HTTP response wrapper. */
export const orvalData = <T>(response: { data: unknown }): T =>
  response.data as T

/** Extracts the items of a list payload from Orval's response wrapper. */
export const orvalItems = <T>(response: { data: unknown }): T[] =>
  itemsOf<T>(response.data)

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
 * Create all fetcher variants every app needs.
 */
export function createAppFetchers(config: AppFetcherConfig): AppFetchers {
  const {
    apiURL,
    authAPIURL = apiURL,
    timeout = 10_000,
  } = config

  const clientConfig = { timeout, fetch: envelopeFetch() }

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

  // An adapter is handed the bare path, so it joins the base URL itself.
  const fetchPublic = envelopeFetch()
  const publicFetcher = createDefaultFetchClient({
    adapter: (path, init) =>
      fetchPublic(`${apiURL.replace(/\/$/, "")}/${path.replace(/^\//, "")}`, init),
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