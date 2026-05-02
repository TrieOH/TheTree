import { AuthInterceptor, type RequestOptions as InterceptorOptions } from "./interceptor";
import type { AuthTokenClaims } from "../utils/token-utils";
import {
  createDefaultFetchClient,
  type DefaultFetchClientConfig,
  type DefaultFetchResult,
  type DefaultSuccessEnvelope,
  type DefaultFailureEnvelope,
  type FetchClient,
  type FetchClientOptions
} from "@soramux/node-fetch-sdk";

export type { DefaultFetchResult as ApiResponse };

export interface ApiRequestOptions extends FetchClientOptions {
  requiresAuth?: boolean;
  skipRefresh?: boolean;
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
    clientConfig?: Omit<DefaultFetchClientConfig, "adapter">,
  ) {
    this.interceptor = new AuthInterceptor({
      baseURL,
      authBaseURL,
      onTokenRefreshed,
    });

    this.client = createDefaultFetchClient({
      ...clientConfig,
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
  clientConfig?: Omit<DefaultFetchClientConfig, "adapter">;
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
  clientConfig?: Omit<DefaultFetchClientConfig, "adapter">;
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
