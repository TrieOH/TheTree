import z from "zod";
import type { VisualMetadata } from "@/shared/ui/MetadataVisualizer";

export const permissionCRUDSchema = z.object({
  id: z.string(),
  project_id: z.string(),
  object: z
    .string()
    .regex(
      /^(\*|[a-zA-Z][a-zA-Z0-9_]*)$/,
      "Object must be '*' or start with a letter and contain only alphanumeric characters or underscores",
    ),
  action: z
    .string()
    .regex(
      /^(\*|[a-zA-Z][a-zA-Z0-9_]*)$/,
      "Action must be '*' or start with a letter and contain only alphanumeric characters or underscores",
    ),
  meta: z.record(z.string(), z.any()).optional(),
});

export type PermissionCRUD = z.infer<typeof permissionCRUDSchema>;

export interface Permission {
  id: string;
  object: string;
  action: string;
  project_id: string;
  meta?: VisualMetadata;
  created_at: string;
}