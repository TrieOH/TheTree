import type {
  AddProjectMemberRequest,
  CreateProjectRequest,
  Project,
  ProjectMember
} from '@trieoh/identityx-models';
import { ProjectRoleAdmin, ProjectRoleMember, ProjectRoleOwner } from '@trieoh/identityx-models';
import { z } from 'zod';

export const projectCreateSchema = z.object({
  name: z.string().min(3, "Project name must be at least 3 characters long"),
  domain: z.url({ error: "Must be a valid URL (e.g., https://example.com)" }),
}) satisfies z.ZodType<CreateProjectRequest>;

export type ProjectCreateI = CreateProjectRequest;

export type ProjectI = Project

// Members

export type ProjectMemberRoleI =
  | typeof ProjectRoleMember
  | typeof ProjectRoleAdmin
  | typeof ProjectRoleOwner;

export const memberAddToProjectSchema = z.object({
  actor_email: z.email({ error: "Must be a valid email address" }),
  role: z.enum([
    ProjectRoleMember,
    ProjectRoleAdmin,
    ProjectRoleOwner
  ], { error: "Invalid role" }),
}) satisfies z.ZodType<AddProjectMemberRequest>;

export type MemberAddToProjectI = AddProjectMemberRequest;

export interface ProjectMemberI
  extends Omit<ProjectMember, "role"> {
  role: ProjectMemberRoleI;
}
