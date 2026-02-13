import z from "zod";

export const versionDraftSchema = z.object({
  schema_id: z.string(),
  project_id: z.string(),
});

export type VersionDraft = z.infer<typeof versionDraftSchema>;

export interface SchemaVersion {
  id: string;
  schema_id: string;
  based_on_version_id: string;
  status: string;
  version_number: number;
  created_at: string;
  updated_at: string;
}