import { createDefaultFetchClient } from "@trieoh/envoy-fetch-ts";
import type {
  ApiKeyCreateResponseI,
  // ApiKeyI,
  CapabilityI,
  ActorI,
  CreateActorRequestI,
  CreateIdentityXAccessClientConfig,
  CreateApiKeyRequest,
} from "./types";
import { env } from "./env";

type FetchClient = ReturnType<typeof createDefaultFetchClient>;

function createClient({ baseURL, apiKey }: CreateIdentityXAccessClientConfig) {
  return createDefaultFetchClient({
    baseURL,
    headers: {
      "X-API-Key": apiKey,
    },
  });
}

export class IdentityXAccessClient {
  private readonly client: FetchClient;

  constructor(config: CreateIdentityXAccessClientConfig) {
    this.client = createClient(config);
  }

  apiKeys = {
    // list: (projectId: string) =>
    //   this.client.get<ApiKeyI[]>(`/projects/${projectId}/api-keys`),
    create: (projectId: string, payload: Omit<CreateApiKeyRequest, 'subject_id'>) =>
      this.client.post<ApiKeyCreateResponseI>(`/projects/${projectId}/api_keys`, payload),
  };

  capabilities = {
    list: (projectId: string) =>
      this.client.get<CapabilityI[]>(`/projects/${projectId}/capabilities`),
  };

  actors = {
    list: (projectId: string) =>
      this.client.get<ActorI[]>(`/projects/${projectId}/actors`),
    getById: (projectId: string, actorId: string) =>
      this.client.get<ActorI>(`/projects/${projectId}/actors/${actorId}`),
    getByEmail: (projectId: string, email: string) =>
      this.client.get<ActorI>(
        `/projects/${projectId}/actors/${encodeURIComponent(email)}:by_email`,
      ),
    create: (projectId: string, payload: CreateActorRequestI) =>
      this.client.post<ActorI>(`/projects/${projectId}/actors`, payload),
  };
}

export function createIdentityXAccessClient(config: CreateIdentityXAccessClientConfig) {
  return new IdentityXAccessClient(config);
}

export const client = new IdentityXAccessClient({
  baseURL: env.BASE_URL,
  apiKey: env.API_KEY,
});
