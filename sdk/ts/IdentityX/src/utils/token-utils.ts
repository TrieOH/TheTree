import type { Api } from "../core/api";

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

  if (typeof window !== "undefined") {
    const expired = "expires=Thu, 01 Jan 1970 00:00:00 GMT";
    document.cookie = `svc_session=; path=/; ${expired}; secure; samesite=lax`;
    document.cookie = `refresh_token=; path=/; ${expired}; secure; samesite=lax`;
  }

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
  const res = await apiInstance.post<{ service_session_id: string, expires_at: string }>(
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
    if (typeof window !== "undefined") {
      const expiresDate = new Date(res.data.expires_at).toUTCString();
      document.cookie = `svc_session=${res.data.service_session_id}; path=/; expires=${expiresDate}; secure; samesite=lax`;
    }

    const claimsRes = await fetchAndSaveClaims(apiInstance, is_up_to_date, true);

    if (typeof window !== "undefined") {
      const refreshExpiry = new Date(claimsRes.data.refresh_expiry_date).toUTCString();
      document.cookie = `refresh_token=${refresh_token}; path=/; expires=${refreshExpiry}; secure; samesite=lax`;
    }

    return claimsRes;
  }

  throw new Error(res.message || "Failed to exchange tokens");
};