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
  profile: ProfileData;
  schema_version: number;
  outdated: boolean;
  updated_at: string;
}

export interface UpsertProfileRequest {
  profile: ProfileData;
}

export interface JsonSchemaProperty {
  type?: "string" | "number" | "integer" | "boolean" | "array" | "object";
  title?: string;
  description?: string;
  format?: string;
  enum?: JsonValue[];
  default?: JsonValue;
  items?: JsonSchemaProperty;
  properties?: Record<string, JsonSchemaProperty>;
  required?: string[];
  [keyword: string]: unknown;
}

export interface ProfileSchema {
  project_id?: string | null;
  schema: JsonSchemaProperty;
  version: number;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface UpsertProfileSchemaRequest {
  schema: JsonSchemaProperty;
  active: boolean;
}
