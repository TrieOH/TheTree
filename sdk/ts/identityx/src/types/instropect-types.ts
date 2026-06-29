export type IntrospectCredentialType = "token" | "api_key";
export type IntrospectActorType = "human" | "service" | "machine";

export interface IntrospectCredential {
  /** Applicable for stateful credentials like api keys */
  id?: string;
  type: IntrospectCredentialType;
}

export interface IntrospectSubject {
  id: string;
  project_id?: string;
  email?: string;
  type: IntrospectActorType;
  capabilities: Record<string, unknown>;
  metadata: Record<string, unknown>;
}

export interface IntrospectResponse {
  cred: IntrospectCredential;
  sub: IntrospectSubject;
}