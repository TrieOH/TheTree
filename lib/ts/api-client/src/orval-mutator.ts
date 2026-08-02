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
}

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

/** Configure the base URLs used by every generated client. */
export const configureApiClient = (next: ApiClientConfig): void => {
  config.baseURL = next.baseURL
  if (next.authBaseURL !== undefined) config.authBaseURL = next.authBaseURL
  api = undefined
}

export interface OrvalRequestOptions {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE" | "HEAD" | "OPTIONS"
  params?: Record<string, unknown>
  body?: string
  headers?: Record<string, string>
  signal?: AbortSignal
}

/** Orval mutator: dispatch a request through the shared fetcher and unwrap the envelope. */
export const customInstance = async <T>(
  url: string,
  options: OrvalRequestOptions,
): Promise<T> => {
  const targetUrl = options.params
    ? `${url}?${new URLSearchParams(
        Object.entries(options.params).flatMap(([key, value]) =>
          Array.isArray(value)
            ? value.map((item) => [key, String(item)])
            : [[key, String(value)]],
        ),
      )}`
    : url

  const result = await getApi().request<unknown>(targetUrl, {
    method: options.method,
    body: options.body,
    headers: options.headers,
    adapterInit: { signal: options.signal },
  })

  if (!result.success) throw new ApiError(result)
  return { data: result.data, status: result.code } as T
}

/** Pass-through body type; generated clients already produce the spec shapes. */
export type BodyType<BodyData> = BodyData

/** Error surfaced by react-query hooks: the shared stack's ApiError. */
export type ErrorType<Error> = ApiError<DefaultFailureEnvelope>
