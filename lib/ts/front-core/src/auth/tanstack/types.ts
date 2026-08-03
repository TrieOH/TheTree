import type { TokenSubject } from "@trieoh/identityx-sdk-ts";

export interface BffIntrospectResponse {
  cred: { id?: string; type: "token" | "api_key" };
  subject: {
    id: string;
    project_id?: string;
    email?: string;
    type: "human" | "service" | "machine";
    capabilities: Record<string, SerializableValue>;
    metadata: Record<string, SerializableValue>;
  };
}

export type SerializableValue =
  | string
  | number
  | boolean
  | null
  | SerializableValue[]
  | { [key: string]: SerializableValue };

export type ProxyHttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
export type IdentityXOAuthProvider = "github" | "google";

export interface IdentityXTransportLogEvent {
  layer: "bff-server" | "bff-client";
  operation: string;
  method: string;
  path: string;
  duration_ms: number;
  success: boolean;
  status?: number;
  error_id?: string;
  message?: string;
}

export interface ServerAuthResult {
  success: boolean;
  code: number;
  message?: string;
  error_id?: string;
  trace?: string[];
  profile?: TokenSubject | null;
}

export interface ServerSessionSnapshot {
  isAuthenticated: boolean;
  profile: TokenSubject | null;
}

export interface ServerProxyRequest {
  path: string;
  target?: "api" | "identityx";
  method?: ProxyHttpMethod;
  body?: SerializableValue;
  headers?: Record<string, string>;
}

export interface ServerOperationResult<T extends SerializableValue = SerializableValue>
  extends ServerAuthResult {
  data?: T;
}

export interface ServerProxyResult<T = SerializableValue> {
  success: boolean;
  code: number;
  data?: T;
  message?: string;
  error_id?: string;
  trace?: string[];
}
