import type {
  ProjectProfileSchema,
  UpsertProfileSchemaRequest,
} from "@trieoh/identityx-api/schemas";

export type ProfileSchemaI = ProjectProfileSchema;
export type ProfileSchemaInputI = UpsertProfileSchemaRequest;

export const DEFAULT_PROFILE_SCHEMA = {
  $schema: "https://json-schema.org/draft/2020-12/schema",
  type: "object",
  properties: {
    full_name: { type: "string", title: "Full name" },
    display_name: { type: "string", title: "Display name" },
    picture_url: { type: "string", format: "uri", title: "Profile picture" },
  },
  required: ["full_name"],
  additionalProperties: false,
} as const;
