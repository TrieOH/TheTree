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

// Stored only in memory
let memoryClaims: AuthTokenClaims | null = null;

export function saveTokenClaims(claims: AuthTokenClaims): void {
  memoryClaims = claims;
  console.log("[TRIEOH SDK] Token claims saved (in-memory)");
}

export function getTokenClaims(): AuthTokenClaims | null {
  return memoryClaims;
}

export function isTokenExpiringSoon(thresholdSeconds: number = 30): boolean {
  const claims = getTokenClaims();
  if (!claims || !claims.access.exp) {
    console.warn("[TRIEOH SDK] No token claims found or exp missing");
    return true;
  }

  const now = Math.floor(Date.now() / 1000);
  const timeUntilExpiry = claims.access.exp - now;
    
  return timeUntilExpiry <= thresholdSeconds;
}

export function isAuthenticated(): boolean {
  const claims = getTokenClaims();
  return claims !== null;
}

export function clearAuthTokens(): void {
  memoryClaims = null;
  console.log("[TRIEOH SDK] Auth tokens and claims cleared");
}

export function getUserInfo() {
  const claims = getTokenClaims();
  if (!claims) return null;
  
  return claims.access.sub
}

export const fetchAndSaveClaims = async (apiInstance: Api) => {
  const res = await apiInstance.get<AuthTokenClaims>("/sessions/me",
    { requiresAuth: true, skipRefresh: true }
  );
  
  if (res.code === 200 && res.data) {
    saveTokenClaims(res.data);
    return res.data;
  }
  clearAuthTokens(); 
  throw new Error(res.message || "Failed to fetch session claims");
};