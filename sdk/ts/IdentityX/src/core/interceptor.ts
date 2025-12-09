import {
  clearAuthTokens,
  isTokenExpiringSoon,
  saveTokenClaims,
  type AuthTokenClaims
} from "../utils/token-utils";
import { env } from "./env";

interface InterceptorConfig {
  baseURL?: string;
  apiKey?: string;
  onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  onRefreshFailed?: (error: Error) => void;
}

export class AuthInterceptor {
  private baseURL: string;
  private apiKey: string;
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;
  private onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  private onRefreshFailed?: (error: Error) => void;

  constructor(config?: InterceptorConfig) {
    this.baseURL = config?.baseURL || env.BASE_URL;
    this.apiKey = config?.apiKey || env.API_KEY;
    this.onTokenRefreshed = config?.onTokenRefreshed;
    this.onRefreshFailed = config?.onRefreshFailed;

    if (!this.apiKey) {
      console.warn("[TRIEOH SDK] API_KEY not found");
      throw new Error("[TRIEOH SDK] API_KEY not found");
    }
  }

  private async refreshToken(): Promise<void> {
    if (this.isRefreshing && this.refreshPromise) return this.refreshPromise;
    
    this.isRefreshing = true;
    this.refreshPromise = (async () => {
      try {
        const response = await fetch(`${this.baseURL}/auth/refresh`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${this.apiKey}`,
          },
          credentials: "include",
        });

        const data = await response.json();

        if (data.code !== 200) {
          if (data.code !== 503) clearAuthTokens();
          throw new Error(data.message || "Failed to refresh token");
        }

        if (data.data) {
          saveTokenClaims(data.data);
          this.onTokenRefreshed?.(data.data);
        }

        console.log("[TRIEOH SDK] Token refreshed successfully");
      } catch (error) {
        console.error("[TRIEOH SDK] Failed to refresh token:", error);
        clearAuthTokens();
        this.onRefreshFailed?.(error as Error);
        throw error;
      } finally {
        this.isRefreshing = false;
        this.refreshPromise = null;
      }
    })();

    return this.refreshPromise;
  }

  async beforeRequest(): Promise<void> {
    if (isTokenExpiringSoon(30)) {
      console.log("[TRIEOH SDK] Token expiring soon, refreshing...");
      try {
        await this.refreshToken();
      } catch (error) {
        console.error("[TRIEOH SDK] Refresh interceptor failed:", error);
      }
    }
  }


  async fetch(url: string, options?: RequestInit): Promise<Response> {
    await this.beforeRequest();
    
    return fetch(url, {
      ...options,
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
    });
  }
}


export const createAuthInterceptor = (config?: InterceptorConfig) => {
  return new AuthInterceptor(config);
};

export const createAuthenticatedFetch = (config?: InterceptorConfig) => {
  const interceptor = new AuthInterceptor(config);
  
  return async (url: string, options?: RequestInit): Promise<Response> => {
    return interceptor.fetch(url, options);
  };
};