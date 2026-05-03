import { authStore } from "../store/auth-store";
import { tokenStore } from "../store/token-store";
import { logger } from "@trieoh/envoy-fetch-ts";
import { browserStorage, cookieStorage } from "./storage-adapter";

export interface AuthTokens {
  AccessTokenString: string;
  RefreshTokenString: string;
  AccessExpiresAt: string;
  RefreshExpiresAt: string;
  Domain: string;
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
let memoryClaims: AuthTokenClaims | null = null;
const ACCESS_EXPIRY_KEY = "trieoh_access_expiry";
const REFRESH_EXPIRY_KEY = "trieoh_refresh_expiry";
const REFRESH_DOMAIN_KEY = "trieoh_refresh_domain";

export function getCookieDomain(returnedDomain?: string) {
  if (typeof window === "undefined") return null;
  const hostname = window.location.hostname;

  // Localhost should never have a domain set for cookies to avoid issues
  if (hostname === 'localhost') return null;

  if (returnedDomain) {
    try {
      let domain = returnedDomain;
      if (domain.startsWith("http")) domain = new URL(domain).hostname;
      return domain;
    } catch {
      return returnedDomain;
    }
  }

  return hostname;
}
export function decodeJwt<T>(token: string): T | null {
  try {
    const payload = token.split(".")[1];
    return JSON.parse(atob(payload));
  } catch {
    return null;
  }
}

export function saveAuthSession(tokens: AuthTokens): void {
  const {
    AccessTokenString,
    RefreshTokenString,
    AccessExpiresAt,
    RefreshExpiresAt,
    Domain
  } = tokens;

  const claims = decodeJwt<TokenClaims>(AccessTokenString);

  if (!claims) {
    logger.error("Failed to decode tokens");
    return;
  }

  tokenStore.setAccessToken(AccessTokenString);

  const refreshExpiry = new Date(RefreshExpiresAt).getTime();
  const accessExpiry = new Date(AccessExpiresAt).getTime();
  const domain = getCookieDomain(Domain);

  cookieStorage.set("refresh_token", RefreshTokenString, {
    expires: new Date(refreshExpiry).toUTCString(),
    domain
  });

  memoryClaims = {
    access_data: claims,
    refresh_expiry_date: refreshExpiry,
  };

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
  if (memoryClaims) return memoryClaims;

  const token = tokenStore.getAccessToken();
  if (!token) return null;

  const claims = decodeJwt<TokenClaims>(token);
  if (!claims) return null;

  // Check if token is expired
  if (claims.exp * 1000 <= Date.now()) return null;

  const refreshExpiryStr = browserStorage.getItem(REFRESH_EXPIRY_KEY);
  if (!refreshExpiryStr) return null;

  memoryClaims = {
    access_data: claims,
    refresh_expiry_date: parseInt(refreshExpiryStr, 10),
  };

  return memoryClaims;
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
  if (!tokenStore.getAccessToken()) return false;
  const expiryStr = browserStorage.getItem(ACCESS_EXPIRY_KEY);
  if (!expiryStr) return false;
  const accessExpiryTimestamp = parseInt(expiryStr, 10);
  return accessExpiryTimestamp > Date.now();
}

export function clearAuthTokens(): void {
  memoryClaims = null;
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
