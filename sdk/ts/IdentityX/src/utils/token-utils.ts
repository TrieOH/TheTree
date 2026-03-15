import type { Api } from "../core/api";
import { authStore } from "../store/auth-store";

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

export function removeCookie(name: string): void {
  setCookie(name, "", "Thu, 01 Jan 1970 00:00:00 GMT");
}

// Stored only in memory
let memoryClaims: AuthTokenClaims | null = null;

export function saveTokenClaims(claims: AuthTokenClaims): void {
  memoryClaims = claims;

  const refreshExpiry = new Date(claims.refresh_expiry_date).getTime();
  const accessExpiry = claims.access_data.exp * 1000;

  if (isNaN(refreshExpiry)) {
    console.error("[TRIEOH SDK] Invalid refresh_expiry_date received:", claims.refresh_expiry_date);
  } else localStorage.setItem(REFRESH_EXPIRY_KEY, String(refreshExpiry));

  if (isNaN(accessExpiry)) {
    console.error("[TRIEOH SDK] Invalid access expiry received:", claims.access_data.exp);
  } else localStorage.setItem(ACCESS_EXPIRY_KEY, String(accessExpiry));

  if (claims.is_up_to_date !== undefined) {
    localStorage.setItem(IS_UP_TO_DATE_KEY, String(claims.is_up_to_date));
  }

  console.log("[TRIEOH SDK] Token claims saved");
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

export function isTokenExpiringSoon(thresholdSeconds: number = 30): boolean {
  try {
    const expiryStr = localStorage.getItem(ACCESS_EXPIRY_KEY);
    if (!expiryStr) return true;

    const accessExpiryTimestamp = parseInt(expiryStr, 10);
    if (isNaN(accessExpiryTimestamp)) return true;

    const now = Date.now();
    const thresholdMs = thresholdSeconds * 1000;

    return (accessExpiryTimestamp - now) <= thresholdMs;
  } catch (e) {
    console.warn("[TRIEOH SDK] Error reading access expiry date:", e);
    return true;
  }
}

export function isRefreshSessionExpired(thresholdSeconds: number = 10): boolean {
  try {
    const expiryStr = localStorage.getItem(REFRESH_EXPIRY_KEY);
    if (!expiryStr) return true;

    const refreshExpiryTimestamp = parseInt(expiryStr, 10);
    if (isNaN(refreshExpiryTimestamp)) return true;

    const now = Date.now();
    const thresholdMs = thresholdSeconds * 1000;

    return (refreshExpiryTimestamp - now) <= thresholdMs;
  } catch (e) {
    console.warn("[TRIEOH SDK] Error reading refresh expiry date:", e);
    return true;
  }
}

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

  console.log("[TRIEOH SDK] Auth tokens and claims cleared");
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
    console.warn("[TRIEOH SDK] fetch claims failed (network/server)", error);
    throw error;
  }
};

export const exchangeAndSaveClaims = async (
  apiInstance: Api,
  access_token: string,
  refresh_token: string,
  is_up_to_date: boolean
) => {
  const res = await apiInstance.post<{ session_id: string, ttl: string }>(
    "/auth/exchange",
    undefined,
    {
      requiresAuth: false,
      headers: {
        "Authorization": `Bearer ${access_token}`,
      },
    }
  );

  if (res.success) {
    const expiresDate = new Date(res.data.ttl).getTime().toString();
    console.log("[TRIEOH SDK] Exchanging tokens, session expires at:", expiresDate);
    setCookie("svc_session", res.data.session_id, expiresDate);

    // FIXME: Apenas o usuário do GoAUTH usa sessions/me
    const claimsRes = await fetchAndSaveClaims(apiInstance, is_up_to_date, true);

    const refreshExpiry = new Date(claimsRes.data.refresh_expiry_date).toUTCString();
    setCookie("refresh_token", refresh_token, refreshExpiry);

    return claimsRes;
  }

  throw new Error(res.message || "Failed to exchange tokens");
};