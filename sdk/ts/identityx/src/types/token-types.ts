export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  access_expires_at: string;
  refresh_expires_at: string;
  domain: string;
}

// FIXME: use this interface for the return type of getUserInfo()
// export type ActorType = "human" | "service" | "machine";

// export interface TokenClaims {
//   sub: {
//     id: string;
//     project_id: string | null;
//     email: string | null;
//     type: ActorType;
//     capabilities: Record<string, unknown> | null;
//     metadata: Record<string, unknown> | null;
//   };
//   iss: string;
//   exp: number;
//   iat: number;
//   jti: string;
// }

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