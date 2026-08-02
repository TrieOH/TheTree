import { OAuthProviderI } from "./common-types";
import { JsonValue } from "./token-types";

export interface OAuthProviderDiscoveryItem {
  id: string;
  provider: OAuthProviderI | string;
  name?: string;
  icon_url?: string;
  enabled?: boolean;
}

export type ProfileData = Record<string, JsonValue>;

export interface ActorProfile {
  actor_id: string;
  project_id?: string;
  profile: ProfileData;
  schema_version?: number;
  created_at: string;
  updated_at: string;
}

export interface UpsertProfileRequest {
  profile: ProfileData;
}