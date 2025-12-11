import {
  clearAuthTokens,
  fetchAndSaveClaims,
  isTokenExpiringSoon,
} from "../utils/token-utils";
import { env } from "./env";

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
  private apiKey: string;
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;

  constructor(baseURL?: string) {
    this.baseURL = baseURL || env.BASE_URL;
    this.apiKey = env.API_KEY;
    if (!this.apiKey) {
      console.warn("[TRIEOH SDK] API_KEY not found, verify your .env file");
      throw new Error("[TRIEOH SDK] API_KEY not found, verify your .env file");
    }
  }

  private get headers() {
    return {
      "Content-Type": "application/json",
      ...(this.apiKey ? { Authorization: `Bearer ${this.apiKey}` } : {}),
    };
  }

  private buildUrl(path: string) {
    return `${this.baseURL.replace(/\/$/, "")}/${path.replace(/^\//, "")}`;
  }

  private async refreshToken(): Promise<void> {    
    if(this.isRefreshing && this.refreshPromise) return this.refreshPromise;
    
    this.isRefreshing = true;
    this.refreshPromise = (async () => {
      try {
        const res = await this.request<string>("/auth/refresh", {
          method: "POST",
          skipRefresh: true,
          requiresAuth: true,
        });
        if (res.code !== 200) {
          if (res.code !== 503) clearAuthTokens();
          throw new Error(res.message || "Failed to refresh token");
        }
        await fetchAndSaveClaims(this);
        console.log("[TRIEOH SDK] Token refreshed successfully");
      } catch (error) {
        console.error("[TRIEOH SDK] Failed to refresh token:", error);
        clearAuthTokens();
        throw error;
      } finally {
        this.isRefreshing = false;
        this.refreshPromise = null;
      }
    })();

    return this.refreshPromise;
  }

  private async interceptRequest(options?: RequestOptions): Promise<void> {
    if (!options?.requiresAuth || options?.skipRefresh) return;

    if (isTokenExpiringSoon(30)) {
      console.log("[TRIEOH SDK] Token expiring soon, refreshing...");
      try {
        await this.refreshToken();
      } catch (error) {
        console.error("[TRIEOH SDK] Refresh interceptor failed:", error);
      }
    }
  }

  async request<T = unknown>(
    path: string,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    try {
      await this.interceptRequest(options);
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
