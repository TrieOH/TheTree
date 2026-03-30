import { z } from 'zod';

export const navigationStateSchema = z.object({
  currentSchemaVersion: z.number().nullable(),
});

export type NavigationStoreState = z.infer<typeof navigationStateSchema>;

export interface NavigationActions {
  setCurrentSchemaVersion: (schemaVersion: number | null) => void;
}
