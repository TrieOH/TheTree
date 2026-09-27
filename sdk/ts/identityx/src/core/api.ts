import { AuthInterceptor, type RequestOptions as InterceptorOptions } from "./interceptor";
import type { AuthTokenClaims } from "../types/token-types";
import type { TokenStore } from "../store/token-store";
import {
  createDefaultFetchClient,
  type DefaultFetchClientConfig,
  type DefaultFetchResult,
  type DefaultSuccessEnvelope,
  type DefaultFailureEnvelope,
  type FetchClient,
  type FetchClientOptions
} from "@trieoh/envoy-fetch-ts";

export type { DefaultFetchResult as ApiResponse };

export interface ApiRequestOptions extends FetchClientOptions {
  requiresAuth?: boolean;
  skipRefresh?: boolean;
}

/**
 * Fetch-client config plus the session storage and transport the auth
 * interceptor uses. Without them the SDK keeps using the browser's storages and
 * the global `fetch`.
 */
export interface ApiClientConfig extends Omit<DefaultFetchClientConfig, "adapter"> {
  tokenStore?: TokenStore;
  fetch?: typeof fetch;
}

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

export class Api {
  readonly interceptor: AuthInterceptor;
  private readonly client: FetchClient<DefaultSuccessEnvelope, DefaultFailureEnvelope>;

  constructor(
    baseURL?: string,
    authBaseURL?: string,
    onTokenRefreshed?: (claims: AuthTokenClaims) => void,
    clientConfig?: ApiClientConfig,
  ) {
    const { tokenStore, fetch: fetchImpl, ...fetchConfig } = clientConfig ?? {};

    this.interceptor = new AuthInterceptor({
      baseURL,
      authBaseURL,
      onTokenRefreshed,
      tokenStore,
      fetch: fetchImpl,
    });

    this.client = createDefaultFetchClient({
      ...fetchConfig,
      adapter: this.interceptor.fetch.bind(this.interceptor),
    });
  }

  request<T>(path: string, options?: ApiRequestOptions) {
    return this.client.request<T>(path, toFetchOptions(options));
  }

  get<T>(path: string, options?: ApiRequestOptions) {
    return this.client.get<T>(path, toFetchOptions(options));
  }

  post<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.post<T>(path, body, toFetchOptions(options));
  }

  put<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.put<T>(path, body, toFetchOptions(options));
  }

  patch<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.patch<T>(path, body, toFetchOptions(options));
  }

  delete<T>(path: string, body?: unknown, options?: ApiRequestOptions) {
    return this.client.delete<T>(path, body, toFetchOptions(options));
  }

  query<T>(path: string, options?: ApiRequestOptions) {
    return this.client.query<T>(path, toFetchOptions(options));
  }
}

export function createFetcher(config?: {
  baseURL?: string;
  authBaseURL?: string;
  clientConfig?: ApiClientConfig;
}) {
  const api = new Api(
    config?.baseURL,
    config?.authBaseURL,
    undefined,
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

export function createQueryFetcher(config?: {
  baseURL?: string;
  authBaseURL?: string;
  clientConfig?: ApiClientConfig;
}) {
  const api = new Api(
    config?.baseURL,
    config?.authBaseURL,
    undefined,
    config?.clientConfig,
  );

  return <TData>(path: string, options?: ApiRequestOptions): Promise<TData> =>
    api.query<TData>(path, options);
}
