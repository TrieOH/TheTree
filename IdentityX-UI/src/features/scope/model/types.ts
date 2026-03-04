import z from "zod";
import type { VisualMetadata } from "@/shared/ui/MetadataVisualizer";

export const scopeCRUDSchema = z.object({
  id: z.string(),
  project_id: z.string(),
  parent_id: z.string().nullable().optional(),
  name: z.string().min(3, "Name must be at least 3 characters long"),
  external_id: z.string().optional()
    .refine(
      val => val === undefined || val === "" || val.length >= 3,
      "External ID must be at least 3 characters long or empty"
    ),
  meta: z.record(z.string(), z.any()).optional(),
});

export type ScopeCRUD = z.infer<typeof scopeCRUDSchema>;

export interface Scope {
  id: string;
  name: string;
  type: string;
  project_id: string;
  parent_id?: string | null;
  external_id?: string;
  meta?: VisualMetadata;
  created_at: string;
}