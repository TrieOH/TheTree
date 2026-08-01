import type { TokenSubject } from "@trieoh/identityx-sdk-ts";

export type SerializableValue =
  | string
  | number
  | boolean
  | null
  | SerializableValue[]
  | { [key: string]: SerializableValue };

export type ProxyHttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
export type IdentityXOAuthProvider = "github" | "google";

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
