import type {
  ApiKey,
  Capability,
  CreateApiKeyRequest,
  CreateApiKeyResponse,
  CreateCapabilityRequest,
} from "@trieoh/identityx-models";

export type { ApiKey, Capability, CreateApiKeyRequest, CreateApiKeyResponse, CreateCapabilityRequest };

export interface IdentityXAccessClientConfig {
  baseURL: string;
  apiKey: string;
}

export interface CreateIdentityXAccessClientConfig extends IdentityXAccessClientConfig {}

export interface ApiKeyCreateResponseI extends CreateApiKeyResponse {}
export interface ApiKeyI extends ApiKey {}
export interface CapabilityI extends Capability {}
export interface CapabilityCreateI extends CreateCapabilityRequest {}
