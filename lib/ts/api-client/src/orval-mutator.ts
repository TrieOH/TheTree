/**
 * Orval custom mutator for generated API clients (lib/ts/<svc>/client).
 *
 * Generated clients call `customInstance(url, { method, params, body })` and
 * receive orval's standard response wrapper `{ data, status }` — where `data`
 * is the unwrapped `fun.Response` payload. Failures reject with the shared
 * stack's `ApiError` (`FetchClientError`, envelope in `.envelope`), so orval
 * clients behave exactly like the hand-written fetchers.
 *
 * Apps must call `configureApiClient` once at startup:
 *
 *   configureApiClient({ baseURL: env.VITE_API_URL, authBaseURL: env.VITE_AUTH_API_URL })
 */
import { ApiError, createFetcher } from "@trieoh/identityx-sdk-ts"
import type { DefaultFailureEnvelope } from "@trieoh/envoy-fetch-ts"

export interface ApiClientConfig {
  /** Base URL for the main API. */
  baseURL: string
  /** Base URL for the authentication API (used for token refresh). */
  authBaseURL?: string
  /** Custom transport, normally used when browser requests go through a BFF. */
  transport?: ApiTransport
  /**
   * Transport for public/unauthenticated operations. Pass the app's
   * `publicFetcher` (via `createOrvalTransport`) here so generated clients
   * called with `{ public: true }` dispatch directly — no auth headers,
   * no BFF proxy.
   */
  publicTransport?: ApiTransport
}

export interface ApiTransportResponse<T = unknown> {
  data: T
  status: number
  headers: Headers
}

export type ApiTransport = <T>(
  url: string,
  options: OrvalRequestOptions,
) => Promise<ApiTransportResponse<T>>

const config: ApiClientConfig = { baseURL: "" }

let api: ReturnType<typeof createFetcher> | undefined

const getApi = (): ReturnType<typeof createFetcher> => {
  if (!api) {
    api = createFetcher({
      baseURL: config.baseURL,
      authBaseURL: config.authBaseURL,
    })
  }
  return api
}

/** Configure the base URLs and transports used by every generated client. */
export const configureApiClient = (next: ApiClientConfig): void => {
  config.baseURL = next.baseURL
  if (next.authBaseURL !== undefined) config.authBaseURL = next.authBaseURL
  config.transport = next.transport
  config.publicTransport = next.publicTransport
  api = undefined
}

export interface OrvalRequestOptions {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE" | "HEAD" | "OPTIONS"
  params?: Record<string, unknown>
  body?: string
  headers?: Record<string, string>
  signal?: AbortSignal
  /**
   * When true, dispatches through the configured `publicTransport` — the
   * app's `publicFetcher` — skipping auth headers and the BFF proxy.
   * Use for public (unauthenticated) routes.
   */
  public?: boolean
}

/** Orval mutator: dispatch a request through the shared fetcher and unwrap the envelope. */
export const customInstance = async <T>(
  url: string,
  options: OrvalRequestOptions,
): Promise<T> => {
  // Public (unauthenticated) routes: skip auth headers and the BFF proxy,
  // dispatching straight to the app's publicFetcher.
  if (options.public) {
    if (config.publicTransport) {
      return config.publicTransport(url, options) as Promise<T>
    }
    throw new Error(
      "customInstance: public option requires configureApiClient({ publicTransport }) — pass the app's publicFetcher via createOrvalTransport.",
    )
  }

  if (config.transport) return config.transport(url, options) as Promise<T>

  const targetUrl = appendParams(url, options.params)

  const result = await getApi().request<unknown>(targetUrl, {
    method: options.method,
    body: options.body,
    headers: options.headers,
    adapterInit: { signal: options.signal },
  })

  if (!result.success) throw new ApiError(result)
  return { data: result.data, status: result.code, headers: new Headers() } as T
}

/** Converts any shared fetch client (direct or BFF-backed) into an Orval transport. */
export const createOrvalTransport = (source: unknown): ApiTransport => {
  const client = source as {
    request(url: string, options: {
      method?: OrvalRequestOptions["method"]
      body?: string
      headers?: Record<string, string>
      adapterInit?: { signal?: AbortSignal }
    }): Promise<{ success: boolean; code: number; data?: unknown }>
  }

  return async <T>(url: string, options: OrvalRequestOptions) => {
    const targetUrl = appendParams(url, options.params)
    const result = await client.request(targetUrl, {
      method: options.method,
      body: options.body,
      headers: options.headers,
      adapterInit: { signal: options.signal },
    })
    if (!result.success) throw new ApiError(result as never)
    return { data: result.data as T, status: result.code, headers: new Headers() }
  }
}

const appendParams = (url: string, params?: Record<string, unknown>): string =>
  params
    ? `${url}?${new URLSearchParams(
      Object.entries(params).flatMap(([key, value]) =>
        Array.isArray(value)
          ? value.map((item) => [key, String(item)])
          : [[key, String(value)]],
      ),
    )}`
    : url

/** Pass-through body type; generated clients already produce the spec shapes. */
export type BodyType<BodyData> = BodyData

/** Error surfaced by react-query hooks: the shared stack's ApiError. */
export type ErrorType<_Error> = ApiError<DefaultFailureEnvelope>