import type { Api, ApiResponse } from "../core/api";
import { AuthTokens } from "../core/services";
import { authStore } from "../store/auth-store";
import { logger } from "@soramux/node-fetch-sdk";
import { browserStorage, cookieStorage } from "./storage-adapter";

export interface TokenClaims {
  sub: {
    id: string;
    email: string;
    session_id: string;
    user_agent: string;
    user_ip: string;
    project_id: string | null;
    verified_at: string | null;
    is_verified: boolean;
    user_type: "client" | "project";
    metadata: Record<string, unknown> | null;
  };
  iss: string;
  exp: number;
  iat: number;
  jti: string;
}

export interface AuthTokenClaims {
  access_data: TokenClaims;
  refresh_expiry_date: string | number;
}

// Stored only in memory
let memoryClaims: AuthTokenClaims | null = null;
const ACCESS_EXPIRY_KEY = "trieoh_access_expiry";
const REFRESH_EXPIRY_KEY = "trieoh_refresh_expiry";

export function getCookieDomain() {
  if (typeof window === "undefined") return null;
  const hostname = window.location.hostname;
  if (hostname === 'localhost') return null;
  if (hostname.endsWith('univents.com.br')) return 'univents.com.br';

  const parts = hostname.split('.');
  if (hostname.endsWith('trieoh.com') && parts.length >= 3) {
    const appName = parts[parts.length - 3];
    return `${appName}.trieoh.com`;
  }

  return hostname;
}

export function saveSession(
  claims: AuthTokenClaims,
  session_id: string,
  session_ttl: string,
  refresh_token: string,
): void {
  cookieStorage.set("svc_session", session_id, {
    expires: new Date(session_ttl).toUTCString(),
    domain: getCookieDomain()
  });

  saveTokenClaims(claims);

  const refreshExpiry = new Date(claims.refresh_expiry_date).toUTCString();
  cookieStorage.set("refresh_token", refresh_token, {
    expires: refreshExpiry,
    domain: getCookieDomain()
  });
}

export const withExchange = async (
  apiInstance: Api,
  res: ApiResponse<AuthTokens>,
  context: "login" | "refresh",
  exchangeURL?: string,
): Promise<ApiResponse<AuthTokens>> => {
  if (!res.success) return res;
  try {
    await exchangeAndSaveClaims(
      apiInstance,
      res.data.access_token,
      res.data.refresh_token,
      exchangeURL,
    );
    return res;
  } catch (error) {
    logger.error(`Exchange failed during ${context}:`, error);
    clearAuthTokens();
    return {
      success: false,
      code: 500,
      message: error instanceof Error ? error.message : "Authentication failed during exchange",
    } as ApiResponse<AuthTokens>;
  }
};

export function saveTokenClaims(claims: AuthTokenClaims): void {
  memoryClaims = claims;

  const refreshExpiry = new Date(claims.refresh_expiry_date).getTime();
  const accessExpiry = claims.access_data.exp * 1000;

  if (isNaN(refreshExpiry)) {
    logger.error("Invalid refresh_expiry_date received:", claims.refresh_expiry_date);
  } else browserStorage.setItem(REFRESH_EXPIRY_KEY, String(refreshExpiry));

  if (isNaN(accessExpiry)) {
    logger.error("Invalid access expiry received:", claims.access_data.exp);
  } else browserStorage.setItem(ACCESS_EXPIRY_KEY, String(accessExpiry));

  authStore.set({
    isAuthenticated: !!claims.access_data,
    isInitializing: false,
  });

  logger.log("Token claims saved");
}

export function getTokenClaims(): AuthTokenClaims | null {
  if (memoryClaims) return memoryClaims;
  return null;
}

function isExpiringSoon(key: string, thresholdSeconds: number): boolean {
  try {
    const expiryStr = browserStorage.getItem(key);
    if (!expiryStr) return true;
    const expiry = parseInt(expiryStr, 10);
    if (isNaN(expiry)) return true;
    return (expiry - Date.now()) <= thresholdSeconds * 1000;
  } catch (e) {
    logger.warn("Error reading expiry:", e);
    return true;
  }
}

export const isTokenExpiringSoon = (t = 30) => isExpiringSoon(ACCESS_EXPIRY_KEY, t);
export const isRefreshSessionExpired = (t = 10) => isExpiringSoon(REFRESH_EXPIRY_KEY, t);

export function isAuthenticated(): boolean {
  const expiryStr = browserStorage.getItem(ACCESS_EXPIRY_KEY);
  if (!expiryStr) return false;
  const accessExpiryTimestamp = parseInt(expiryStr, 10);
  return accessExpiryTimestamp > Date.now();
}

export function clearAuthTokens(): void {
  memoryClaims = null;
  browserStorage.removeItem(ACCESS_EXPIRY_KEY);
  browserStorage.removeItem(REFRESH_EXPIRY_KEY);

  const domain = getCookieDomain();
  cookieStorage.remove("svc_session", domain);
  cookieStorage.remove("refresh_token", domain);

  authStore.reset();

  logger.log("Auth tokens and claims cleared");
}

export function getUserInfo() {
  const claims = getTokenClaims();
  if (!claims) return null;

  return claims.access_data.sub
}

export const fetchAndSaveClaims = async (
  apiInstance: Api,
  skipRefresh?: boolean
) => {
  try {
    const res = await apiInstance.get<AuthTokenClaims>("/sessions/me",
      { requiresAuth: true, skipRefresh }
    );
    if (res.success) {
      saveTokenClaims(res.data);
      return { ...res, data: res.data };
    }
    throw new Error(res.message || "Failed to fetch session claims");
  } catch (error) {
    logger.warn("fetch claims failed (network/server)", error);
    throw error;
  }
};

export function decodeJwtExp(token: string): number | null {
  try {
    const payload = token.split(".")[1];
    const decoded = JSON.parse(atob(payload));
    return decoded.exp ?? null;
  } catch {
    return null;
  }
}

export const exchangeAndSaveClaims = async (
  apiInstance: Api,
  access_token: string,
  refresh_token: string,
  exchangeURL?: string
) => {
  // Project: Custom URL
  if (exchangeURL) {
    const res = await apiInstance.post<{ session_id: string, ttl: string, claims: TokenClaims }>(
      exchangeURL,
      undefined,
      {
        requiresAuth: false,
        headers: { "Authorization": `Bearer ${access_token}` },
      }
    );

    if (res.success) {
      const refreshExp = decodeJwtExp(refresh_token);
      const claims: AuthTokenClaims = {
        access_data: res.data.claims,
        refresh_expiry_date: refreshExp ? refreshExp * 1000 : 0,
      };

      saveSession(claims, res.data.session_id, res.data.ttl, refresh_token);

      return { ...res, data: claims };
    }

    throw new Error(res.message || "Failed to exchange tokens (project)");
  }

  // Default: /auth/exchange + sessions/me
  const res = await apiInstance.post<{ session_id: string, ttl: string }>(
    "/auth/exchange",
    undefined,
    {
      requiresAuth: false,
      headers: { "Authorization": `Bearer ${access_token}` },
    }
  );

  if (res.success) {
    cookieStorage.set("svc_session", res.data.session_id, {
      expires: new Date(res.data.ttl).toUTCString(),
      domain: getCookieDomain()
    });

    const claimsRes = await fetchAndSaveClaims(apiInstance, true);

    const refreshExpiry = new Date(claimsRes.data.refresh_expiry_date).toUTCString();
    cookieStorage.set("refresh_token", refresh_token, {
      expires: refreshExpiry,
      domain: getCookieDomain()
    });

    return claimsRes;
  }

  throw new Error(res.message || "Failed to exchange tokens");
};
