import { z } from 'zod';
import type { JsonValue } from '@/shared/model/types';

export const projectCRUDSchema = z.object({
  id: z.string(),
  project_name: z.string().min(3, "Project name must be at least 3 characters long"),
  domain: z.url({error: "Invalid URL format"}),
});

export type ProjectCRUD = z.infer<typeof projectCRUDSchema>;


export interface Project {
  id: string;
  project_name: string;
  domain: string;
  is_active: boolean;
  owner_id: string;
  metadata: Record<string, JsonValue>;
  created_at: string;
  updated_at: string;
}