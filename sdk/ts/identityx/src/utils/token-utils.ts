import { authStore } from "../store/auth-store";
import { defaultTokenStore, type TokenStore } from "../store/token-store";
import { logger } from "@trieoh/envoy-fetch-ts";
import type {
  AuthTokenClaims,
  AuthTokens,
  TokenClaims,
  TokenSubject,
} from "../types/token-types";

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

export function saveAuthSession(
  tokens: AuthTokens,
  store: TokenStore = defaultTokenStore,
): void {
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

  const accessExpiry = new Date(access_expires_at).getTime();
  const refreshExpiry = new Date(refresh_expires_at).getTime();

  store.setAccessToken(access_token);
  store.setRefreshToken(refresh_token);
  store.setAccessExpiry(accessExpiry);
  store.setRefreshExpiry(refreshExpiry);
  store.setClaims({
    access_data: claims,
    refresh_expiry_date: refreshExpiry,
  });

  authStore.set({
    isAuthenticated: true,
    isInitializing: false,
  });

  logger.log("Auth session saved");
}

export function getStoredRefreshToken(
  store: TokenStore = defaultTokenStore,
): string | null {
  return store.getRefreshToken();
}

export function getTokenClaims(
  store: TokenStore = defaultTokenStore,
): AuthTokenClaims | null {
  const cached = store.getClaims();
  if (cached) return cached;

  const token = store.getAccessToken();
  if (!token) return null;

  const claims = decodeJwt<TokenClaims>(token);
  if (!claims) return null;

  // Check if token is expired
  if (claims.exp * 1000 <= Date.now()) return null;

  const refreshExpiry = store.getRefreshExpiry();
  if (refreshExpiry === null) return null;

  const sessionData: AuthTokenClaims = {
    access_data: claims,
    refresh_expiry_date: refreshExpiry,
  };

  store.setClaims(sessionData);

  return sessionData;
}

function isExpiringSoon(
  readExpiry: () => number | null,
  thresholdSeconds: number,
): boolean {
  try {
    const expiry = readExpiry();
    if (expiry === null) return true;
    return (expiry - Date.now()) <= thresholdSeconds * 1000;
  } catch (e) {
    logger.warn("Error reading expiry:", e);
    return true;
  }
}

export const isTokenExpiringSoon = (t = 30, store: TokenStore = defaultTokenStore) =>
  isExpiringSoon(() => store.getAccessExpiry(), t);
export const isRefreshSessionExpired = (t = 10, store: TokenStore = defaultTokenStore) =>
  isExpiringSoon(() => store.getRefreshExpiry(), t);

export function isAuthenticated(store: TokenStore = defaultTokenStore): boolean {
  if (!store.getAccessToken()) return false;
  const accessExpiry = store.getAccessExpiry();
  if (accessExpiry === null) return false;
  return accessExpiry > Date.now();
}

export function clearAuthTokens(store: TokenStore = defaultTokenStore): void {
  store.clear();

  authStore.reset();

  logger.log("Auth tokens and claims cleared");
}

export function getUserInfo(
  store: TokenStore = defaultTokenStore,
): TokenSubject | null {
  const claims = getTokenClaims(store);
  if (!claims) return null;

  return claims.access_data.subject;
}

export function decodeJwtExp(token: string): number | null {
  const decoded = decodeJwt<{ exp: number }>(token);
  return decoded?.exp ?? null;
}
