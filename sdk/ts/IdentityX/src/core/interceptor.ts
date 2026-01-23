import {
  clearAuthTokens,
  isRefreshSessionExpired,
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
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;
  private onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  private onRefreshFailed?: (error: Error) => void;

  constructor(config?: InterceptorConfig) {
    this.baseURL = config?.baseURL || env.BASE_URL;
    this.onTokenRefreshed = config?.onTokenRefreshed;
    this.onRefreshFailed = config?.onRefreshFailed;
  }

  private async fetchClaimsAndSave(): Promise<AuthTokenClaims> {
    const response = await fetch(`${this.baseURL}/sessions/me`, {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    });

    const resJson = await response.json();

    if (resJson.code !== 200 || !resJson.data) {
      clearAuthTokens();
      throw new Error(resJson.message || "Failed to fetch session claims after refresh");
    }
    
    const claims = resJson.data as AuthTokenClaims;
    saveTokenClaims(claims);
    return claims;
  }

  private async refreshToken(): Promise<void> {
    if (this.isRefreshing && this.refreshPromise) return this.refreshPromise;
    
    this.isRefreshing = true;
    this.refreshPromise = (async () => {
      let claims: AuthTokenClaims;
      try {
        const response = await fetch(`${this.baseURL}/auth/refresh`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });

        const data = await response.json();

        if (data.code !== 200) {
          if (data.code !== 503) clearAuthTokens();
          throw new Error(data.message || "Failed to refresh token");
        }

        claims = await this.fetchClaimsAndSave();
        this.onTokenRefreshed?.(claims);
        console.log("[TRIEOH SDK] Token refreshed successfully");
      } catch (error) {
        console.warn("[TRIEOH SDK] Failed to refresh token:", error);
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
    if(isRefreshSessionExpired()) {
      clearAuthTokens();
      return;
    }
    
    if (isTokenExpiringSoon(30)) {
      console.log("[TRIEOH SDK] Token expiring soon, refreshing...");
      try {
        await this.refreshToken();
      } catch (error) {
        console.warn("[TRIEOH SDK] Refresh interceptor failed:", error);
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