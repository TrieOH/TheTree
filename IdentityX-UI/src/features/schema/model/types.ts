import { z } from 'zod';

export const schemaCRUDSchema = z.object({
  id: z.string(),
  project_id: z.string(),
  title: z.string().min(3, "Title must be at least 3 characters long"),
  flow_id: z.string().min(3, "Flow ID must be at least 3 characters long")
});

export type SchemaCRUD = z.infer<typeof schemaCRUDSchema>;

export interface Schema {
  id: string;
  title: string;
  flow_id: string;
  status: string;
  type: string;
  current_version_id: string;
  project_id: string;
  created_at: string;
  updated_at: string;
}
