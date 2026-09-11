import { joinUrl } from "../utils/url-utils";
import {
  clearAuthTokens,
  isRefreshSessionExpired,
  isTokenExpiringSoon,
  saveAuthSession,
  getTokenClaims,
  getStoredRefreshToken,
} from "../utils/token-utils";
import { env } from "./env";
import { logger } from "@trieoh/envoy-fetch-ts";
import { defaultTokenStore, type TokenStore } from "../store/token-store";
import type { AuthTokenClaims, AuthTokens } from "../types/token-types";

export interface RequestOptions extends RequestInit {
  requiresAuth?: boolean;
  skipRefresh?: boolean;
}

interface InterceptorConfig {
  baseURL?: string;
  authBaseURL?: string;
  onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  onRefreshFailed?: (error: Error) => void;
  tokenStore?: TokenStore;
  /**
   * Transport for every IdentityX call, refresh included. Inject one to add
   * tracing, retries or a runtime binding instead of patching `globalThis.fetch`.
   */
  fetch?: typeof fetch;
}

interface RefreshEnvelope {
  code: number;
  data?: AuthTokens;
  message?: string;
}

async function readRefreshEnvelope(response: Response): Promise<RefreshEnvelope> {
  try {
    const body = (await response.json()) as Partial<RefreshEnvelope>;
    return {
      code: typeof body.code === "number" ? body.code : response.status,
      ...(body.data ? { data: body.data } : {}),
      ...(body.message ? { message: body.message } : {}),
    };
  } catch {
    return { code: response.status, message: response.statusText };
  }
}

export class AuthInterceptor {
  private baseURL: string;
  private authBaseURL: string;
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;
  private onTokenRefreshed?: (claims: AuthTokenClaims) => void;
  private onRefreshFailed?: (error: Error) => void;
  /** Session this interceptor reads and writes. */
  readonly tokenStore: TokenStore;
  private readonly fetchImpl: typeof fetch;

  constructor(config?: InterceptorConfig) {
    this.baseURL = config?.baseURL || env.BASE_URL;
    this.authBaseURL = config?.authBaseURL || this.baseURL;
    this.onTokenRefreshed = config?.onTokenRefreshed;
    this.onRefreshFailed = config?.onRefreshFailed;
    this.tokenStore = config?.tokenStore ?? defaultTokenStore;
    this.fetchImpl = config?.fetch ?? fetch;
  }

  async refreshToken(): Promise<void> {
    if (this.isRefreshing && this.refreshPromise) return this.refreshPromise;

    this.isRefreshing = true;
    this.refreshPromise = (async () => {
      let shouldClear = false;
      try {
        const refreshToken = getStoredRefreshToken(this.tokenStore);
        if (!refreshToken) {
          shouldClear = true;
          throw new Error("No refresh token available");
        }

        const response = await this.fetchImpl(
          joinUrl(this.authBaseURL, "/auth/refresh"),
          {
            method: "POST",
            credentials: "omit",
            headers: {
              "Content-Type": "application/json",
              "Refresh-Token": refreshToken,
            },
          }
        );
        const res = await readRefreshEnvelope(response);
        const isSuccessfulCode = res.code >= 200 && res.code < 300;
        if (!isSuccessfulCode || !res.data || !res.data.access_token) {
          shouldClear = res.code >= 400 && res.code < 500;
          throw new Error(res.message || "Failed to refresh token");
        }

        saveAuthSession(res.data, this.tokenStore);

        const claims = getTokenClaims(this.tokenStore);
        if (claims) this.onTokenRefreshed?.(claims);

        logger.log("Token refreshed successfully");
      } catch (error) {
        logger.warn("Failed to refresh token:", error);
        if (shouldClear) clearAuthTokens(this.tokenStore);
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
    if (isRefreshSessionExpired(10, this.tokenStore)) {
      clearAuthTokens(this.tokenStore);
      return;
    }

    const hasAccessToken = !!this.tokenStore.getAccessToken();

    if (!hasAccessToken || isTokenExpiringSoon(30, this.tokenStore)) {
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
      const accessToken = this.tokenStore.getAccessToken();
      const headers = new Headers(options?.headers);

      if (shouldAuth && accessToken) {
        headers.set("Authorization", `Bearer ${accessToken}`);
      }
      if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
      }

      return this.fetchImpl(finalUrl, {
        ...options,
        headers,
        credentials: "omit",
      });
    };

    let response = await executeFetch();

    if (response.status === 401 && shouldAuth && !isRefreshReq) {
      const hasRefreshToken = !!getStoredRefreshToken(this.tokenStore);

      if (hasRefreshToken) {
        logger.log("401 detected, attempting token refresh...");
        try {
          await this.refreshToken();
          response = await executeFetch();
        } catch (e) {
          logger.error("Retry failed after refresh error");
        }
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
