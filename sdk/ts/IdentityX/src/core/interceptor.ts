import { joinUrl } from "../utils/url-utils";
import {
  clearAuthTokens,
  isRefreshSessionExpired,
  isTokenExpiringSoon,
  saveAuthSession,
  getTokenClaims,
  type AuthTokenClaims,
  type AuthTokens
} from "../utils/token-utils";
import { env } from "./env";
import { logger, simpleFetch } from "@soramux/node-fetch-sdk";
import { cookieStorage } from "../utils/storage-adapter";
import { tokenStore } from "../store/token-store";

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
    this.baseURL = config?.baseURL || env.BASE_URL;
    this.authBaseURL = config?.authBaseURL || this.baseURL;
    this.onTokenRefreshed = config?.onTokenRefreshed;
    this.onRefreshFailed = config?.onRefreshFailed;
  }

  private async refreshToken(): Promise<void> {
    if (this.isRefreshing && this.refreshPromise) return this.refreshPromise;

    this.isRefreshing = true;
    this.refreshPromise = (async () => {
      try {
        const res = await simpleFetch<{ code: number; data?: AuthTokens; message?: string }>(
          joinUrl(this.authBaseURL, "/auth/refresh"),
          { method: "POST", credentials: "include" }
        );

        if (res.code !== 200 || !res.data) {
          if (res.code !== 503) clearAuthTokens();
          throw new Error(res.message || "Failed to refresh token");
        }

        const { access_token, refresh_token } = res.data;
        saveAuthSession(access_token, refresh_token);

        const claims = getTokenClaims();
        if (claims) this.onTokenRefreshed?.(claims);

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

    const hasAccessToken = !!tokenStore.getAccessToken();
    const hasRefreshToken = !!cookieStorage.get("refresh_token");

    if ((!hasAccessToken && hasRefreshToken) || (hasAccessToken && isTokenExpiringSoon(30))) {
      try {
        await this.refreshToken();
      } catch (error) {
        logger.warn("Proactive refresh failed:", error);
      }
    }
  }

  async fetch(url: string, options?: RequestOptions): Promise<Response> {
    const shouldAuth = options?.requiresAuth !== false;
    const isRefreshReq = url.includes("/auth/refresh");

    if (shouldAuth && !isRefreshReq && !options?.skipRefresh) {
      await this.beforeRequest();
    }

    const finalUrl = joinUrl(this.baseURL, url);

    const executeFetch = async (): Promise<Response> => {
      const accessToken = tokenStore.getAccessToken();
      const headers = new Headers(options?.headers);

      if (shouldAuth && accessToken) {
        headers.set("Authorization", `Bearer ${accessToken}`);
      }
      if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
      }

      return fetch(finalUrl, {
        ...options,
        headers,
        credentials: "include",
      });
    };

    let response = await executeFetch();

    if (response.status === 401 && shouldAuth && !isRefreshReq) {
      logger.log("401 detected, attempting one-time retry after refresh...");
      try {
        await this.refreshToken();
        response = await executeFetch();
      } catch (e) {
        logger.error("Retry failed after refresh error");
      }
    }

    return response;
  }
}

export const createAuthInterceptor = (config?: InterceptorConfig) => new AuthInterceptor(config);

export const createAuthenticatedFetch = (config?: InterceptorConfig) => {
  const interceptor = new AuthInterceptor(config);
  return (url: string, options?: RequestOptions) => interceptor.fetch(url, options);
};
