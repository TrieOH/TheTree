import { joinUrl } from "../utils/url-utils";
import {
  clearAuthTokens,
  decodeJwtExp,
  isRefreshSessionExpired,
  isTokenExpiringSoon,
  saveSession,
  saveTokenClaims,
  setCookie,
  TokenClaims,
  type AuthTokenClaims
} from "../utils/token-utils";
import { env } from "./env";
import { logger } from "@soramux/node-fetch-sdk";
import { simpleFetch } from "@soramux/node-fetch-sdk";


export interface RequestOptions extends RequestInit {
  requiresAuth?: boolean;
  skipRefresh?: boolean;
}

interface InterceptorConfig {
  baseURL?: string;
  authBaseURL?: string;
  exchangeURL?: string;
  onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  onRefreshFailed?: (error: Error) => void;
}

export class AuthInterceptor {
  private baseURL: string;
  private authBaseURL: string;
  private exchangeURL?: string;
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;
  private onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  private onRefreshFailed?: (error: Error) => void;

  constructor(config?: InterceptorConfig) {
    this.baseURL = config?.baseURL || env.BASE_URL;
    this.authBaseURL = config?.authBaseURL || this.baseURL;
    this.exchangeURL = config?.exchangeURL;
    this.onTokenRefreshed = config?.onTokenRefreshed;
    this.onRefreshFailed = config?.onRefreshFailed;
  }

  private async fetchClaimsAndSave(is_up_to_date?: boolean): Promise<AuthTokenClaims> {
    const res = await simpleFetch<{ data?: AuthTokenClaims, code: number }>(
      joinUrl(this.authBaseURL, "/sessions/me")
    );

    if (res.code !== 200 || !res.data) {
      clearAuthTokens();
      throw new Error("Failed to fetch session claims after refresh");
    }

    const claims = { ...res.data, is_up_to_date };
    saveTokenClaims(claims);
    return claims;
  }

  private async exchange(
    access_token: string,
    refresh_token: string,
    is_up_to_date: boolean
  ): Promise<AuthTokenClaims> {
    // Project: Custom URL
    if (this.exchangeURL) {
      const res = await simpleFetch<{
        data?: { session_id: string; ttl: string; claims: TokenClaims },
        code: number,
        message?: string
      }>(
        this.exchangeURL,
        { method: "POST", headers: { "Authorization": `Bearer ${access_token}` } }
      );

      if (res.code !== 200 || !res.data) {
        clearAuthTokens();
        throw new Error(res.message || "Failed to exchange tokens (project)");
      }

      const refreshExp = decodeJwtExp(refresh_token);
      const claims: AuthTokenClaims = {
        access_data: res.data.claims,
        is_up_to_date,
        refresh_expiry_date: refreshExp ? refreshExp * 1000 : 0,
      };

      saveSession(claims, res.data.session_id, res.data.ttl, refresh_token);
      return claims;
    }

    // Default: /auth/exchange + sessions/me
    const res = await simpleFetch<{ code: number; message?: string; data?: { session_id: string; expires_at: string } }>(
      joinUrl(this.authBaseURL, "/auth/exchange"),
      { method: "POST", headers: { "Authorization": `Bearer ${access_token}` } }
    );

    if (res.code !== 200 || !res.data) {
      clearAuthTokens();
      throw new Error(res.message || "Failed to exchange tokens");
    }

    const expiresDate = new Date(res.data.expires_at).getTime().toString();
    setCookie("svc_session", res.data.session_id, expiresDate);

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
        logger.log("Token refreshed successfully");
      } catch (error) {
        logger.warn("Failed to refresh token:", error);
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
      logger.log("Token expiring soon, refreshing...");
      try {
        await this.refreshToken();
      } catch (error) {
        logger.warn("Refresh interceptor failed:", error);
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