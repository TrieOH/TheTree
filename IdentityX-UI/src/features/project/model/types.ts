// import { z } from 'zod';

// export const projectSchema = z.object({
//   id: z.string(),
//   name: z.string().min(3, 'Nome deve ter pelo menos 3 caracteres'),
// });

// export type ProjectFormData = z.infer<typeof projectSchema>;

// export const projectDefaultValues: ProjectFormData = {
//   id: '',
//   name: '',
// };

export interface Project {
  id: string;
  project_name: string;
  is_active: boolean;
  owner_id: string;
  // I need to include metadata
  created_at: string;
  updated_at: string;
}