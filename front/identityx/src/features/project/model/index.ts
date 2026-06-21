import type { CreateProjectRequest, Project } from '@trieoh/identityx-models';
import { z } from 'zod';

export const projectCreateSchema = z.object({
  name: z.string().min(3, "Project name must be at least 3 characters long"),
  domain: z.url({ error: "Must be a valid URL (e.g., https://example.com)" }),
}) satisfies z.ZodType<CreateProjectRequest>;

export type ProjectCreateI = CreateProjectRequest;

export type ProjectI = Project