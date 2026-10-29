import { createDefaultFetchClient, type DefaultFetchResult } from "@trieoh/envoy-fetch-ts";
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
import { itemsOf, toFailure, toSuccess } from "./wire";

type FetchClient = ReturnType<typeof createDefaultFetchClient>;

function createClient({ baseURL, apiKey, fetch: fetchImpl }: CreateIdentityXAccessClientConfig) {
  return createDefaultFetchClient({
    baseURL,
    toSuccess,
    toFailure,
    // An adapter is handed the bare path, so it joins the base URL itself.
    ...(fetchImpl
      ? {
          adapter: (path: string, init?: RequestInit) =>
            fetchImpl(`${baseURL.replace(/\/$/, "")}/${path.replace(/^\//, "")}`, init),
        }
      : {}),
    headers: {
      "X-API-Key": apiKey,
    },
  });
}

// Callers get the array; IdentityX sends a list as { items }.
function items<T>(res: DefaultFetchResult<unknown>): DefaultFetchResult<T[]> {
  return res.success ? { ...res, data: itemsOf<T>(res.data) } : res;
}

export class IdentityXAccessClient {
  private readonly client: FetchClient;

  constructor(config: CreateIdentityXAccessClientConfig) {
    this.client = createClient(config);
  }

  apiKeys = {
    // list: (projectId: string) =>
    //   this.client.get<ApiKeyI[]>(`/projects/${projectId}/api-keys`),
    create: (projectId: string, payload: CreateApiKeyRequest) =>
      this.client.post<ApiKeyCreateResponseI>(`/projects/${projectId}/api_keys`, payload),
  };

  capabilities = {
    list: (projectId: string) =>
      this.client.get<unknown>(`/projects/${projectId}/capabilities`).then(items<CapabilityI>),
  };

  actors = {
    list: (projectId: string) =>
      this.client.get<unknown>(`/projects/${projectId}/actors`).then(items<ActorI>),
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
