import { env } from "./env";
import { AuthInterceptor, type RequestOptions } from "./interceptor";
import type { AuthTokenClaims } from "../utils/token-utils";

/**
 * Base structure shared by all API responses.
 */
export interface BaseResponse {
  module: string;
  message: string;
  timestamp: string;
  code: number;
}

/**
 * Standardized API Response. 
 * Use 'success' to narrow down to the data or error details.
 */
export type ApiResponse<T> = 
  | (BaseResponse & { success: true; data: T })
  | (BaseResponse & { 
      success: false; 
      error_id: string; 
      trace?: string[]; 
    });

/**
 * Internal type for parsing the raw response from the server.
 */
type RawApiResponse<T> = BaseResponse & { 
  data?: T; 
  error_id?: string; 
  trace?: string[] 
};

/**
 * Custom error class for API failures.
 * Provides access to the full standardized response.
 */
export class ApiError extends Error {
  code: number;
  trace?: string[];
  response: ApiResponse<unknown>;

  constructor(response: ApiResponse<unknown>) {
    super(response.message);
    this.name = "ApiError";
    this.code = response.code;
    this.trace = !response.success ? response.trace : [];
    this.response = response;
  }
}

export class Api {
  private baseURL: string;
  private authInterceptor: AuthInterceptor;

  constructor(baseURL?: string, authBaseURL?: string, onTokenRefreshed?: (claims: AuthTokenClaims) => void) {
    this.baseURL = baseURL || env.BASE_URL;
    this.authInterceptor = new AuthInterceptor({ 
      baseURL: this.baseURL,
      authBaseURL: authBaseURL,
      onTokenRefreshed
    });
  }

  private get headers() {
    return { "Content-Type": "application/json" };
  }

  async request<T = unknown>(path: string, options?: RequestOptions): Promise<ApiResponse<T>> {
    try {
      const response = await this.authInterceptor.fetch(path, {
        ...options,
        headers: { ...this.headers, ...(options?.headers ?? {}) },
      });

      const raw: RawApiResponse<T> = await response.json().catch(() => ({
        module: "Client",
        message: response.statusText || "Unknown error",
        timestamp: new Date().toISOString(),
        code: response.status,
        error_id: "CLIENT_PARSE_ERROR",
        trace: ["Failed to parse API response as JSON"],
      }));

      if (!response.ok) {
        return {
          success: false,
          module: raw.module || "Unknown",
          message: raw.message || response.statusText || "An unknown error occurred",
          timestamp: raw.timestamp || new Date().toISOString(),
          code: raw.code || response.status,
          error_id: raw.error_id || "UNKNOWN_ERROR",
          trace: raw.trace,
        };
      }

      return {
        success: true,
        module: raw.module,
        message: raw.message,
        timestamp: raw.timestamp,
        code: raw.code,
        data: raw.data as T,
      };
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "A network or unknown error occurred.";
      return {
        success: false,
        module: "Network",
        message: errorMessage,
        timestamp: new Date().toISOString(),
        code: 503,
        error_id: "CLIENT_NETWORK_ERROR",
        trace: error instanceof Error ? [error.stack || errorMessage] : [errorMessage],
      };
    }
  }

  get<T = unknown>(path: string, options?: RequestOptions) {
    return this.request<T>(path, { ...options, method: "GET" });
  }

  post<T = unknown>(path: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>(path, {
      ...options,
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  put<T = unknown>(path: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>(path, {
      ...options,
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  patch<T = unknown>(path: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>(path, {
      ...options,
      method: "PATCH",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  delete<T = unknown>(path: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>(path, { 
      ...options, 
      method: "DELETE", 
      body: body ? JSON.stringify(body) : undefined 
    });
  }
}

/**
 * Creates a fetcher object with all HTTP methods returning standardized ApiResponse.
 * 
 * @param config - Optional configuration for the fetcher (baseURL, authBaseURL, etc.)
 * @returns An object with get, post, put, patch, delete, and request methods.
 */
export const createFetcher = (config?: { baseURL?: string; authBaseURL?: string }) => {
  const api = new Api(config?.baseURL, config?.authBaseURL);

  return {
    request: api.request.bind(api),
    get: api.get.bind(api),
    post: api.post.bind(api),
    put: api.put.bind(api),
    patch: api.patch.bind(api),
    delete: api.delete.bind(api),
  };
};

/**
 * Creates a simple fetcher function that returns data directly or throws an ApiError on failure.
 * Ideal for Query libraries like TanStack Query.
 * 
 * @param config - Optional configuration for the fetcher (baseURL, authBaseURL, etc.)
 * @returns A single fetcher function that returns a Promise of TData.
 */
export const createQueryFetcher = (config?: { baseURL?: string; authBaseURL?: string }) => {
  const api = new Api(config?.baseURL, config?.authBaseURL);

  return async <TData>(path: string, init?: RequestOptions): Promise<TData> => {
    const response = await api.request<TData>(path, init);
    if (!response.success) throw new ApiError(response);
    return response.data;
  };
};
