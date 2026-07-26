export interface ApiKey {
  id: string;
  subject_id: string;
  name: string;
  display_prefix: string;
  key_hash: string;
  metadata: unknown;
  expires_at?: string;
  revoked_at?: string;
  last_used_at?: string;
  created_by: string;
  created_at: string;
}

export interface Actor {
  id: string;
  project_id?: string;
  auth_method: string;
  verified_at?: string;
  email?: string;
  type: string;
  metadata?: unknown;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface Capability {
  id: string;
  project_id?: string;
  resource: string;
  action: string;
  created_by: string;
  created_at: string;
}

export interface CreateApiKeyRequest {
  subject_id?: string;
  capabilities: string[];
  name: string;
  env: string;
  expires_at?: string;
}

export interface CreateApiKeyResponse {
  key?: ApiKey;
  raw_key: string;
}

export interface CreateCapabilityRequest {
  resource: string;
  action: string;
}

export interface CreateActorRequest {
  auth_method: string;
  type: string;
  email?: string;
}

export interface IdentityXAccessClientConfig {
  baseURL: string;
  apiKey: string;
}

export interface CreateIdentityXAccessClientConfig extends IdentityXAccessClientConfig {}

export interface ApiKeyCreateResponseI extends CreateApiKeyResponse { }
export interface ApiKeyI extends ApiKey { }
export interface CapabilityI extends Capability { }
export interface CapabilityCreateI extends CreateCapabilityRequest { }
export interface ActorI extends Actor { }
export interface CreateActorRequestI extends CreateActorRequest { }
