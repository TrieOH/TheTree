import { env } from "./env";
import { AuthInterceptor } from "./interceptor";

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  module: string;
  timestamp: string;
  trace?: string[];
  data?: T;
}

interface RequestOptions extends RequestInit {
  requiresAuth?: boolean;
  skipRefresh?: boolean;
}

export class Api {
  private baseURL: string;
  private authInterceptor: AuthInterceptor;

  constructor(baseURL?: string) {
    this.baseURL = baseURL || env.BASE_URL;
    this.authInterceptor = new AuthInterceptor({ baseURL: this.baseURL });
  }

  private get headers() {
    return { "Content-Type": "application/json" };
  }

  private buildUrl(path: string) {
    return `${this.baseURL.replace(/\/$/, "")}/${path.replace(/^\//, "")}`;
  }

  async request<T = unknown>(
    path: string,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    try {
      if(options?.requiresAuth && !options.skipRefresh) 
        await this.authInterceptor.beforeRequest();
      const res = await fetch(this.buildUrl(path), {
        ...options,
        headers: { ...this.headers, ...(options?.headers ?? {}) },
        credentials: "include"
      });

      const data = await res.json();

      return data as ApiResponse<T>;
    } catch (error) {
      return {
        code: 503,
        message: "Network request failed — API may be offline.",
        module: "network",
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

  delete<T = unknown>(path: string, options?: RequestOptions) {
    return this.request<T>(path, { ...options, method: "DELETE" });
  }
}
