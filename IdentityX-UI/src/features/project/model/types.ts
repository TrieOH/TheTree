import { z } from 'zod';

export const projectCRUDSchema = z.object({
  id: z.string(),
  project_name: z.string().min(3, "Project name must be at least 3 characters long"),
});

export type ProjectCRUD = z.infer<typeof projectCRUDSchema>;


export interface Project {
  id: string;
  project_name: string;
  is_active: boolean;
  owner_id: string;
  // I need to include metadata
  created_at: string;
  updated_at: string;
}