interface TokenClaims {
  sub: {
    id: string;
    email: string;
  };
  iss: string;
  exp: number;
  iat: number;
  jti: string;
}

interface RefreshTokenClaims {
  sub: {
    access_jti: string;
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
  access_token_claims: TokenClaims;
  refresh_token_claims: RefreshTokenClaims;
}

const TOKEN_CLAIMS_KEY = "auth_token_claims";

export function saveTokenClaims(claims: AuthTokenClaims): void {
  if (typeof window === "undefined") return;
  
  try {
    localStorage.setItem(TOKEN_CLAIMS_KEY, JSON.stringify(claims));
    console.log("[TRIEOH SDK] Token claims saved");
  } catch (error) {
    console.error("[TRIEOH SDK] Failed to save token claims:", error);
  }
}

export function getTokenClaims(): AuthTokenClaims | null {
  if (typeof window === "undefined") return null;
  
  try {
    const claims = localStorage.getItem(TOKEN_CLAIMS_KEY);
    if (!claims) return null;
    return JSON.parse(claims) as AuthTokenClaims;
  } catch (error) {
    console.error("[TRIEOH SDK] Failed to get token claims:", error);
    return null;
  }
}

export function isTokenExpiringSoon(thresholdSeconds: number = 30): boolean {
  const claims = getTokenClaims();
  if (!claims || !claims.access_token_claims.exp) {
    console.warn("[TRIEOH SDK] No token claims found or exp missing");
    return true;
  }

  const now = Math.floor(Date.now() / 1000);
  const timeUntilExpiry = claims.access_token_claims.exp - now;
  
  console.log(`[TRIEOH SDK] Token expires in ${timeUntilExpiry}s`);
  
  return timeUntilExpiry <= thresholdSeconds;
}

export function isAuthenticated(): boolean {
  const claims = getTokenClaims();
  return claims !== null;
}

export function clearAuthTokens(): void {
  if (typeof window !== "undefined") localStorage.removeItem(TOKEN_CLAIMS_KEY);
  console.log("[TRIEOH SDK] Auth tokens and claims cleared");
}

export function getUserInfo(): { id: string; email: string } | null {
  const claims = getTokenClaims();
  if (!claims) return null;
  
  return {
    id: claims.access_token_claims.sub.id,
    email: claims.access_token_claims.sub.email,
  };
}