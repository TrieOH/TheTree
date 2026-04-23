import { authStore } from "../store/auth-store";
import { tokenStore } from "../store/token-store";
import { logger } from "@soramux/node-fetch-sdk";
import { browserStorage, cookieStorage } from "./storage-adapter";

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
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

export function getCookieDomain() {
  if (typeof window === "undefined") return null;
  const hostname = window.location.hostname;
  if (hostname === 'localhost') return null;
  if (hostname.endsWith('univents.com.br')) return 'univents.com.br';
  if (hostname.endsWith('identityx.com.br')) return 'identityx.com.br';

  const parts = hostname.split('.');
  if (hostname.endsWith('trieoh.com') && parts.length >= 3) {
    const appName = parts[parts.length - 3];
    return `${appName}.trieoh.com`;
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

export function saveAuthSession(
  access_token: string,
  refresh_token: string,
): void {
  const claims = decodeJwt<TokenClaims>(access_token);
  const refreshClaims = decodeJwt<{ exp: number }>(refresh_token);

  if (!claims || !refreshClaims) {
    logger.error("Failed to decode tokens");
    return;
  }

  tokenStore.setAccessToken(access_token);

  const refreshExpiry = refreshClaims.exp * 1000;
  const accessExpiry = claims.exp * 1000;

  cookieStorage.set("refresh_token", refresh_token, {
    expires: new Date(refreshExpiry).toUTCString(),
    domain: getCookieDomain()
  });

  memoryClaims = {
    access_data: claims,
    refresh_expiry_date: refreshExpiry,
  };

  browserStorage.setItem(ACCESS_EXPIRY_KEY, String(accessExpiry));
  browserStorage.setItem(REFRESH_EXPIRY_KEY, String(refreshExpiry));

  authStore.set({
    isAuthenticated: true,
    isInitializing: false,
  });

  logger.log("Auth session saved");
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

  const domain = getCookieDomain();
  cookieStorage.remove("refresh_token", domain);

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
