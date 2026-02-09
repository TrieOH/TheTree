import { z } from 'zod';

export const navigationStateSchema = z.object({
  currentProjectId: z.string().nullable(),
});

export type NavigationStoreState = z.infer<typeof navigationStateSchema>;

export interface NavigationActions {
  setCurrentProjectId: (projectId: string | null) => void;
}
