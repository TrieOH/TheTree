import { AuthInterceptor, type RequestOptions as InterceptorOptions } from "./interceptor";
import type { AuthTokenClaims } from "../utils/token-utils";
import {
  createDefaultFetchClient,
  type DefaultFetchClientConfig,
  type DefaultFetchResult,
  type DefaultSuccessEnvelope,
  type DefaultFailureEnvelope,
  type FetchClient,
  type FetchClientOptions,
  type createFetchClient,
  type FetchClientError as ApiError
} from "@soramux/node-fetch-sdk";
// Re-export public surface under the original names so existing consumers
// don't need to update their imports.
export type { DefaultFetchResult as ApiResponse };

// ─── ApiRequestOptions ────────────────────────────────────────────────────────

/**
 * Per-request options for the {@link Api} class and its factory helpers.
 *
 * Extends the base {@link FetchClientOptions} with auth-layer controls that
 * are forwarded to the {@link AuthInterceptor} via `adapterInit`.
 * `fetch-utils` itself has no knowledge of these fields.
 */
export interface ApiRequestOptions extends FetchClientOptions {
  /**
   * Set to `false` to bypass auth handling entirely for this request.
   * The interceptor will skip the token-refresh check and not attach
   * credentials. Defaults to `true`.
   */
  requiresAuth?: boolean;
  /**
   * Set to `true` to skip the proactive token-refresh check for this request
   * even when auth is required. Useful for endpoints that are always fast
   * and tolerate a slightly stale token.
   */
  skipRefresh?: boolean;
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

/**
 * Merges auth-specific fields from {@link ApiRequestOptions} into
 * `adapterInit` so they are forwarded to the interceptor without polluting
 * the core `FetchClientOptions` interface.
 */
function toFetchOptions(options?: ApiRequestOptions): FetchClientOptions | undefined {
  if (!options) return undefined;

  const { requiresAuth, skipRefresh, ...rest } = options;

  const interceptorFields: Partial<InterceptorOptions> = {};
  if (requiresAuth !== undefined) interceptorFields.requiresAuth = requiresAuth;
  if (skipRefresh !== undefined) interceptorFields.skipRefresh = skipRefresh;

  return {
    ...rest,
    adapterInit: {
      ...rest.adapterInit,
      ...interceptorFields,
    },
  };
}

// ─── Api ──────────────────────────────────────────────────────────────────────

/**
 * Thin façade that wires an {@link AuthInterceptor} to a
 * {@link FetchClient} built with the default envelope scheme.
 *
 * **Responsibilities of this class:**
 * - Constructing and owning the `AuthInterceptor` instance.
 * - Bridging auth-specific request options (`requiresAuth`, `skipRefresh`)
 *   into the adapter call via `adapterInit`.
 * - Exposing a familiar `get / post / put / patch / delete / request / query`
 *   surface typed with {@link ApiRequestOptions}.
 *
 * All response parsing, timeout handling, and error shaping are handled by
 * `fetch-utils`. For custom envelope shapes use {@link createFetchClient}
 * from `fetch-utils` directly.
 *
 * For most use-cases prefer the factory helpers {@link createFetcher} or
 * {@link createQueryFetcher} over instantiating this class directly.
 */
export class Api {
  /**
   * The underlying auth interceptor.
   * Exposed for advanced use (e.g. calling `interceptor.beforeRequest`
   * manually or reusing it in a custom fetch client).
   */
  readonly interceptor: AuthInterceptor;

  private readonly client: FetchClient<DefaultSuccessEnvelope, DefaultFailureEnvelope>;

  /**
   * @param baseURL           - Base URL for all API requests.
   * @param authBaseURL       - Base URL for auth endpoints (`/auth/refresh`,
   *   `/auth/exchange`, `/sessions/me`). Defaults to `baseURL` when omitted.
   * @param onTokenRefreshed  - Callback invoked whenever a token refresh
   *   succeeds. Receives the fresh `AuthTokenClaims`.
   * @param exchangeURL       - Custom token-exchange URL. When set the
   *   interceptor uses project-mode exchange instead of `/auth/exchange`.
   * @param clientConfig      - Extra config forwarded to
   *   {@link createDefaultFetchClient} (e.g. `timeout`, custom adapters).
   */
  constructor(
    baseURL?: string,
    authBaseURL?: string,
    onTokenRefreshed?: (claims: AuthTokenClaims) => void,
    exchangeURL?: string,
    clientConfig?: Omit<DefaultFetchClientConfig, "adapter">,
  ) {
    this.interceptor = new AuthInterceptor({
      baseURL,
      authBaseURL,
      onTokenRefreshed,
      exchangeURL,
    });

    this.client = createDefaultFetchClient({
      ...clientConfig,
      // The interceptor manages baseURL internally — don't pass it again.
      adapter: this.interceptor.fetch.bind(this.interceptor),
    });
  }

  /** @see {@link FetchClient.request} */
  request<T>(path: string, options?: ApiRequestOptions) {
    return this.client.request<T>(path, toFetchOptions(options));
  }

  /** @see {@link FetchClient.get} */
  get<T>(path: string, options?: ApiRequestOptions) {
    return this.client.get<T>(path, toFetchOptions(options));
  }

  /** @see {@link FetchClient.post} */
  post<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.post<T>(path, body, toFetchOptions(options));
  }

  /** @see {@link FetchClient.put} */
  put<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.put<T>(path, body, toFetchOptions(options));
  }

  /** @see {@link FetchClient.patch} */
  patch<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.patch<T>(path, body, toFetchOptions(options));
  }

  /** @see {@link FetchClient.delete} */
  delete<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.delete<T>(path, body, toFetchOptions(options));
  }

  /**
   * Like {@link request}, but resolves to `TData` directly and throws
   * {@link ApiError} on failure.
   *
   * @see {@link FetchClient.query}
   */
  query<T>(path: string, options?: ApiRequestOptions) {
    return this.client.query<T>(path, toFetchOptions(options));
  }
}

// ─── Factory helpers ──────────────────────────────────────────────────────────

/**
 * Creates an object exposing all HTTP methods and returning
 * {@link DefaultFetchResult} discriminated unions. Good for call-sites that
 * prefer to handle success and failure in the same control flow.
 *
 * ```ts
 * const api = createFetcher({ baseURL: "https://api.example.com" });
 *
 * const result = await api.get<User>("/users/me");
 * if (result.success) {
 *   console.log(result.data.name);
 * } else {
 *   console.error(result.error_id);
 * }
 * ```
 *
 * @param config - Optional base URLs, custom exchange URL, and client config.
 */
export function createFetcher(config?: {
  baseURL?: string;
  authBaseURL?: string;
  exchangeURL?: string;
  clientConfig?: Omit<DefaultFetchClientConfig, "adapter">;
}) {
  const api = new Api(
    config?.baseURL,
    config?.authBaseURL,
    undefined,
    config?.exchangeURL,
    config?.clientConfig,
  );

  return {
    request: api.request.bind(api),
    get: api.get.bind(api),
    post: api.post.bind(api),
    put: api.put.bind(api),
    patch: api.patch.bind(api),
    delete: api.delete.bind(api),
    query: api.query.bind(api),
  };
}

/**
 * Creates a single async function that resolves to `TData` on success and
 * throws a {@link ApiError} on failure.
 *
 * Designed for query libraries such as **TanStack Query**, where the
 * `queryFn` is expected to throw on failure rather than return a
 * discriminated union.
 *
 * ```ts
 * const fetcher = createQueryFetcher({ baseURL: "https://api.example.com" });
 *
 * // Inside a TanStack Query definition:
 * useQuery({
 *   queryKey: ["user"],
 *   queryFn: () => fetcher<User>("/users/me"),
 * });
 * ```
 *
 * @param config - Optional base URLs and custom exchange URL.
 */
export function createQueryFetcher(config?: {
  baseURL?: string;
  authBaseURL?: string;
  exchangeURL?: string;
  clientConfig?: Omit<DefaultFetchClientConfig, "adapter">;
}) {
  const api = new Api(
    config?.baseURL,
    config?.authBaseURL,
    undefined,
    config?.exchangeURL,
    config?.clientConfig,
  );

  return <TData>(path: string, options?: ApiRequestOptions): Promise<TData> =>
    api.query<TData>(path, options);
}