import type { Api } from "../core/api";

interface TokenClaims {
  sub: {
    id: string;
    email: string;
    session_id: string;
    user_agent: string;
    user_ip: string;
  };
  iss: string;
  exp: number;
  iat: number;
  jti: string;
}

export interface AuthTokenClaims {
  access: TokenClaims;
  refresh_expire_date: number;
}

const ACCESS_EXPIRY_KEY = "trieoh_access_expiry";
const REFRESH_EXPIRY_KEY = "trieoh_refresh_expiry";

// Stored only in memory
let memoryClaims: AuthTokenClaims | null = null;

export function saveTokenClaims(claims: AuthTokenClaims): void {
  memoryClaims = claims;
  localStorage.setItem(REFRESH_EXPIRY_KEY, String(claims.refresh_expire_date));
  localStorage.setItem(ACCESS_EXPIRY_KEY, String(claims.access.exp));
  console.log("[TRIEOH SDK] Token claims saved");
}

export function getTokenClaims(): AuthTokenClaims | null {
  return memoryClaims;
}

export function isTokenExpiringSoon(thresholdSeconds: number = 30): boolean {
  try {
    const expiryStr = localStorage.getItem(ACCESS_EXPIRY_KEY);
    if (!expiryStr) return true;

    const accessExpiryTimestamp = parseInt(expiryStr, 10);
    const now = Math.floor(Date.now() / 1000);
    const timeUntilExpiry = accessExpiryTimestamp - now;
    
    return timeUntilExpiry <= thresholdSeconds;
  } catch (e) {
    console.warn("[TRIEOH SDK] Error reading access expiry date:", e);
    return true;
  }
}

export function isRefreshSessionExpired(): boolean {
  try {
    const expiryStr = localStorage.getItem(REFRESH_EXPIRY_KEY);
    if (!expiryStr) return true;

    const refreshExpiryTimestamp = parseInt(expiryStr, 10);
    const now = Math.floor(Date.now() / 1000);
    
    return refreshExpiryTimestamp <= now;
  } catch (e) {
    console.warn("[TRIEOH SDK] Error reading refresh expiry date:", e);
    return true;
  }
}

export function isAuthenticated(): boolean {
  const claims = getTokenClaims();
  if (!claims) return false;
  const now = Math.floor(Date.now() / 1000);
  return claims.access.exp > now;
}

export function clearAuthTokens(): void {
  memoryClaims = null;
  localStorage.removeItem(ACCESS_EXPIRY_KEY);
  localStorage.removeItem(REFRESH_EXPIRY_KEY);
  console.log("[TRIEOH SDK] Auth tokens and claims cleared");
}

export function getUserInfo() {
  const claims = getTokenClaims();
  if (!claims) return null;
  
  return claims.access.sub
}

export const fetchAndSaveClaims = async (apiInstance: Api) => {
  try {
    const res = await apiInstance.get<AuthTokenClaims>("/sessions/me",
      { requiresAuth: true }
    );
    
    if (res.code === 200 && res.data) {
      saveTokenClaims(res.data);
      return { code: 200 };
    }
    throw new Error(res.message || "Failed to fetch session claims");
  } catch (error) {
    console.warn("[TRIEOH SDK] fetch claims failed (network/server)", error);
    throw error;
  }
};