import { z } from 'zod';

export const navigationStateSchema = z.object({
  currentProjectId: z.string().nullable(),
  currentSchemaId: z.string().nullable(),
  currentSchemaVersion: z.number().nullable(),
});

export type NavigationStoreState = z.infer<typeof navigationStateSchema>;

export interface NavigationActions {
  setCurrentProjectId: (projectId: string | null) => void;
  setCurrentSchemaId: (schemaId: string | null) => void;
  setCurrentSchemaVersion: (schemaVersion: number | null) => void;
}
