import type { Api, ApiResponse } from "../core/api";
import { AuthTokens } from "../core/services";
import { authStore } from "../store/auth-store";
import { logger } from "@soramux/node-fetch-sdk";

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
  is_up_to_date?: boolean;
}

// Stored only in memory
let memoryClaims: AuthTokenClaims | null = null;
const ACCESS_EXPIRY_KEY = "trieoh_access_expiry";
const REFRESH_EXPIRY_KEY = "trieoh_refresh_expiry";
const IS_UP_TO_DATE_KEY = "trieoh_is_up_to_date";

function getCookieDomain(hostname: string) {
  if (hostname.endsWith('univents.com.br')) return 'univents.com.br';

  const parts = hostname.split('.');
  if (hostname.endsWith('trieoh.com') && parts.length >= 3) {
    const appName = parts[parts.length - 3];
    return `${appName}.trieoh.com`;
  }

  return hostname;
}

export function setCookie(name: string, value: string, expires?: string) {
  if (typeof window === "undefined") return;

  const hostname = window.location.hostname;
  const isSecure = window.location.protocol === 'https:';
  const domain = hostname !== 'localhost' ? getCookieDomain(hostname) : null;

  const cookieParts = [
    `${name}=${value}`,
    domain ? `Domain=${domain}` : '',
    `Path=/`,
    isSecure ? 'SameSite=None' : 'SameSite=Lax',
    isSecure ? 'Secure' : '',
    expires ? `expires=${expires}` : '',
  ];

  const cookieString = cookieParts.filter(Boolean).join('; ');
  document.cookie = cookieString;
}

export function getCookie(name: string): string | null {
  if (typeof window === "undefined") return null;
  const nameEQ = name + "=";
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

export function removeCookie(name: string): void {
  setCookie(name, "", "Thu, 01 Jan 1970 00:00:00 GMT");
}

export function saveSession(
  claims: AuthTokenClaims,
  session_id: string,
  session_ttl: string,
  refresh_token: string,
): void {
  const sessionExpiry = new Date(session_ttl).getTime().toString();
  setCookie("svc_session", session_id, sessionExpiry);
  saveTokenClaims(claims);
  const refreshExpiry = new Date(claims.refresh_expiry_date).toUTCString();
  setCookie("refresh_token", refresh_token, refreshExpiry);
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
      res.data.is_up_to_date,
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
  } else localStorage.setItem(REFRESH_EXPIRY_KEY, String(refreshExpiry));

  if (isNaN(accessExpiry)) {
    logger.error("Invalid access expiry received:", claims.access_data.exp);
  } else localStorage.setItem(ACCESS_EXPIRY_KEY, String(accessExpiry));

  if (claims.is_up_to_date !== undefined) {
    localStorage.setItem(IS_UP_TO_DATE_KEY, String(claims.is_up_to_date));
  }

  authStore.set({
    isAuthenticated: !!claims.access_data,
    isUpToDate: claims.is_up_to_date ?? isUpToDate(),
  });

  logger.log("Token claims saved");
}

export function getTokenClaims(): AuthTokenClaims | null {
  if (memoryClaims) return memoryClaims;
  return null;
}

export function isUpToDate(): boolean {
  if (memoryClaims?.is_up_to_date !== undefined) return memoryClaims.is_up_to_date;
  const stored = localStorage.getItem(IS_UP_TO_DATE_KEY);
  return stored === "true";
}

function isExpiringSoon(key: string, thresholdSeconds: number): boolean {
  try {
    const expiryStr = localStorage.getItem(key);
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
  const expiryStr = localStorage.getItem(ACCESS_EXPIRY_KEY);
  if (!expiryStr) return false;
  const accessExpiryTimestamp = parseInt(expiryStr, 10);
  return accessExpiryTimestamp > Date.now();
}

export function clearAuthTokens(): void {
  memoryClaims = null;
  localStorage.removeItem(ACCESS_EXPIRY_KEY);
  localStorage.removeItem(REFRESH_EXPIRY_KEY);
  localStorage.removeItem(IS_UP_TO_DATE_KEY);

  removeCookie("svc_session");
  removeCookie("refresh_token");

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
  is_up_to_date?: boolean,
  skipRefresh?: boolean
) => {
  try {
    const res = await apiInstance.get<AuthTokenClaims>("/sessions/me",
      { requiresAuth: true, skipRefresh }
    );
    if (res.success) {
      const claims = { ...res.data, is_up_to_date: is_up_to_date ?? isUpToDate() };
      saveTokenClaims(claims);
      return { ...res, data: claims };
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
  is_up_to_date: boolean,
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
        is_up_to_date,
        refresh_expiry_date: refreshExp ? refreshExp * 1000 : 0,
      };

      const expiresDate = new Date(res.data.ttl).getTime().toString();
      setCookie("svc_session", res.data.session_id, expiresDate);

      saveTokenClaims(claims);

      const refreshExpiry = new Date(claims.refresh_expiry_date).toUTCString();
      setCookie("refresh_token", refresh_token, refreshExpiry);

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
    const expiresDate = new Date(res.data.ttl).getTime().toString();
    setCookie("svc_session", res.data.session_id, expiresDate);

    const claimsRes = await fetchAndSaveClaims(apiInstance, is_up_to_date, true);

    const refreshExpiry = new Date(claimsRes.data.refresh_expiry_date).toUTCString();
    setCookie("refresh_token", refresh_token, refreshExpiry);

    return claimsRes;
  }

  throw new Error(res.message || "Failed to exchange tokens");
};
