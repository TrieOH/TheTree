export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  access_expires_at: string;
  refresh_expires_at: string;
  domain: string;
}

export type ActorType = "human" | "service" | "machine";

export interface TokenSubject {
  id: string;
  project_id: string | null;
  email: string | null;
  type: ActorType;
  capabilities: Record<string, unknown> | null;
  metadata: Record<string, unknown> | null;
}

export interface TokenClaims {
  subject: TokenSubject;
  iss: string;
  exp: number;
  iat: number;
  jti: string;
}

export interface AuthTokenClaims {
  access_data: TokenClaims;
  refresh_expiry_date: string | number;
}
