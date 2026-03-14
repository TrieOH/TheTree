import { joinUrl } from "../utils/url-utils";
import {
  clearAuthTokens,
  isRefreshSessionExpired,
  isTokenExpiringSoon,
  saveTokenClaims,
  setCookie,
  type AuthTokenClaims
} from "../utils/token-utils";
import { env } from "./env";


export interface RequestOptions extends RequestInit {
  requiresAuth?: boolean;
  skipRefresh?: boolean;
}

interface InterceptorConfig {
  baseURL?: string;
  authBaseURL?: string;
  onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  onRefreshFailed?: (error: Error) => void;
}

export class AuthInterceptor {
  private baseURL: string;
  private authBaseURL: string;
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;
  private onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  private onRefreshFailed?: (error: Error) => void;

  constructor(config?: InterceptorConfig) {
    this.authBaseURL = config?.authBaseURL || env.BASE_URL;
    this.baseURL = config?.baseURL || this.authBaseURL;

    this.onTokenRefreshed = config?.onTokenRefreshed;
    this.onRefreshFailed = config?.onRefreshFailed;
  }

  private async fetchClaimsAndSave(is_up_to_date?: boolean): Promise<AuthTokenClaims> {
    const response = await fetch(joinUrl(this.authBaseURL, "/sessions/me"), {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    });

    const resJson = await response.json();

    if (resJson.code !== 200 || !resJson.data) {
      clearAuthTokens();
      throw new Error(resJson.message || "Failed to fetch session claims after refresh");
    }

    const claims = { ...resJson.data, is_up_to_date } as AuthTokenClaims;
    saveTokenClaims(claims);
    return claims;
  }

  private async exchange(access_token: string, refresh_token: string, is_up_to_date: boolean): Promise<AuthTokenClaims> {
    const response = await fetch(joinUrl(this.authBaseURL, "/auth/exchange"), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${access_token}`,
      },
      credentials: "include",
    });

    const res = await response.json();

    if (res.code !== 200) {
      clearAuthTokens();
      throw new Error(res.message || "Failed to exchange tokens for session");
    }

    const expiresDate = new Date(res.data.expires_at).toUTCString();
    setCookie("svc_session", res.data.service_session_id, expiresDate);

    const claims = await this.fetchClaimsAndSave(is_up_to_date);

    const refreshExpiry = new Date(claims.refresh_expiry_date).toUTCString();
    setCookie("refresh_token", refresh_token, refreshExpiry);

    return claims;
  }

  private async refreshToken(): Promise<void> {
    if (this.isRefreshing && this.refreshPromise) return this.refreshPromise;

    this.isRefreshing = true;
    this.refreshPromise = (async () => {
      try {
        const response = await fetch(joinUrl(this.authBaseURL, "/auth/refresh"), {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });

        const data = await response.json();

        if (data.code !== 200) {
          if (data.code !== 503) clearAuthTokens();
          throw new Error(data.message || "Failed to refresh token");
        }

        const { access_token, refresh_token, is_up_to_date } = data.data;

        const claims = await this.exchange(access_token, refresh_token, is_up_to_date);

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
    if (isRefreshSessionExpired()) {
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

  async fetch(url: string, options?: RequestOptions): Promise<Response> {
    const shouldAuth = options?.requiresAuth !== false;
    const shouldRefresh = !options?.skipRefresh;

    if (shouldAuth && shouldRefresh) await this.beforeRequest();

    const finalUrl = joinUrl(this.baseURL, url);

    return fetch(finalUrl, {
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

  return async (url: string, options?: RequestOptions): Promise<Response> => {
    return interceptor.fetch(url, options);
  };
};