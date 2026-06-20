import { authStore } from "../store/auth-store";
import { tokenStore } from "../store/token-store";
import { logger } from "@trieoh/envoy-fetch-ts";
import { browserStorage, cookieStorage } from "./storage-adapter";

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  access_expires_at: string;
  refresh_expires_at: string;
  domain: string;
}

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
let _cachedClaims: AuthTokenClaims | null = null;
const ACCESS_EXPIRY_KEY = "trieoh_access_expiry";
const REFRESH_EXPIRY_KEY = "trieoh_refresh_expiry";
const REFRESH_DOMAIN_KEY = "trieoh_refresh_domain";

export function getCookieDomain(returnedDomain?: string) {
  if (typeof window === "undefined") return null;
  const hostname = window.location.hostname;
  const isLocalhost = hostname === "localhost" || hostname === "127.0.0.1" || hostname.includes("localhost");
  if (isLocalhost) return null;

  if (returnedDomain) {
    try {
      let domain = returnedDomain;
      if (domain.startsWith("http")) domain = new URL(domain).hostname;
      if (hostname === domain || hostname.endsWith("." + domain)) return domain;
      return hostname;
    } catch { return hostname; }
  }

  return hostname;
}

export function decodeJwt<T>(token: string): T | null {
  try {
    const payload = token.split(".")[1];
    if (!payload) return null;

    let base64 = payload.replace(/-/g, "+").replace(/_/g, "/");
    while (base64.length % 4) {
      base64 += "=";
    }

    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );

    return JSON.parse(jsonPayload);
  } catch (error) {
    logger.error("Error decoding JWT:", error);
    return null;
  }
}

export function saveAuthSession(tokens: AuthTokens): void {
  const {
    access_token,
    refresh_token,
    access_expires_at,
    refresh_expires_at,
  } = tokens;

  const claims = decodeJwt<TokenClaims>(access_token);

  if (!claims) {
    logger.error("Failed to decode tokens");
    return;
  }

  tokenStore.setAccessToken(access_token);

  const refreshExpiry = new Date(refresh_expires_at).getTime();
  const accessExpiry = new Date(access_expires_at).getTime();
  const domain = getCookieDomain(tokens.domain);

  cookieStorage.set("refresh_token", refresh_token, {
    expires: new Date(refreshExpiry).toUTCString(),
    domain
  });

  const sessionData: AuthTokenClaims = {
    access_data: claims,
    refresh_expiry_date: refreshExpiry,
  };

  _cachedClaims = sessionData;

  browserStorage.setItem(ACCESS_EXPIRY_KEY, String(accessExpiry));
  browserStorage.setItem(REFRESH_EXPIRY_KEY, String(refreshExpiry));
  if (domain) browserStorage.setItem(REFRESH_DOMAIN_KEY, domain);
  else browserStorage.removeItem(REFRESH_DOMAIN_KEY);

  authStore.set({
    isAuthenticated: true,
    isInitializing: false,
  });

  logger.log("Auth session saved");
}

export function getTokenClaims(): AuthTokenClaims | null {
  if (_cachedClaims) return _cachedClaims;

  const token = tokenStore.getAccessToken();
  if (!token) return null;

  const claims = decodeJwt<TokenClaims>(token);
  if (!claims) return null;

  // Check if token is expired
  if (claims.exp * 1000 <= Date.now()) return null;

  const refreshExpiryStr = browserStorage.getItem(REFRESH_EXPIRY_KEY);
  if (!refreshExpiryStr) return null;

  const refreshExpiry = Number(refreshExpiryStr);
  if (isNaN(refreshExpiry)) return null;

  const sessionData = {
    access_data: claims,
    refresh_expiry_date: refreshExpiry,
  };

  _cachedClaims = sessionData;

  return sessionData;
}

function isExpiringSoon(key: string, thresholdSeconds: number): boolean {
  try {
    const stored = browserStorage.getItem(key);
    if (!stored) return true;
    const expiry = Number(stored);
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
  if (!tokenStore.getAccessToken()) return false;
  const stored = browserStorage.getItem(ACCESS_EXPIRY_KEY);
  if (!stored) return false;
  const accessExpiryTimestamp = Number(stored);
  if (isNaN(accessExpiryTimestamp)) return false;
  return accessExpiryTimestamp > Date.now();
}

export function clearAuthTokens(): void {
  _cachedClaims = null;
  tokenStore.clear();
  browserStorage.removeItem(ACCESS_EXPIRY_KEY);
  browserStorage.removeItem(REFRESH_EXPIRY_KEY);

  const domain = browserStorage.getItem(REFRESH_DOMAIN_KEY) || getCookieDomain();
  cookieStorage.remove("refresh_token", domain);
  browserStorage.removeItem(REFRESH_DOMAIN_KEY);

  authStore.reset();

  logger.log("Auth tokens and claims cleared");
}

export function getUserInfo() {
  const claims = getTokenClaims();
  if (!claims) return null;

  return claims.access_data.sub
}

export function decodeJwtExp(token: string): number | null {
  const decoded = decodeJwt<{ exp: number }>(token);
  return decoded?.exp ?? null;
}
