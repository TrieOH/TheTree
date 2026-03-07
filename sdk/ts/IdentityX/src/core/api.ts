import { env } from "./env";
import { AuthInterceptor, type RequestOptions } from "./interceptor";

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  module: string;
  timestamp: string;
  trace?: string[];
  data?: T;
}

export class Api {
  private baseURL: string;
  private authInterceptor: AuthInterceptor;

  constructor(baseURL?: string, authBaseURL?: string) {
    this.baseURL = baseURL || env.BASE_URL;
    this.authInterceptor = new AuthInterceptor({ 
      baseURL: this.baseURL,
      authBaseURL: authBaseURL
    });
  }

  private get headers() {
    return { "Content-Type": "application/json" };
  }

  async request<T = unknown>(
    path: string,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    try {
      const res = await this.authInterceptor.fetch(path, {
        ...options,
        headers: { ...this.headers, ...(options?.headers ?? {}) },
      });

      const data = await res.json();
      return data as ApiResponse<T>;
    } catch (error) {
      return {
        code: 503,
        message: (error as Error).message || "Network request failed — API may be offline.",
        module: "Network",
        timestamp: new Date().toISOString(),
        trace: [(error as Error).message || "Unknown network error"],
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
