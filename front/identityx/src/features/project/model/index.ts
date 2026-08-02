import type {
  AddProjectMemberRequest,
  CreateProjectRequest,
  Project,
  ProjectMember,
} from "@trieoh/identityx-api/schemas";
import { ProjectRole } from "@trieoh/identityx-api/schemas";
import { z } from "zod";

export const projectCreateSchema = z.object({
  name: z.string().min(3, "Project name must be at least 3 characters long"),
  domain: z.url({ error: "Must be a valid URL (e.g., https://example.com)" }),
  brand_slug: z
    .string()
    .min(3, "Brand slug must be at least 3 characters long"),
}) satisfies z.ZodType<CreateProjectRequest>;

export type ProjectCreateI = CreateProjectRequest;

export type ProjectI = Project;

// Members

export type ProjectMemberRoleI = "member" | "admin" | "owner";

export const memberAddToProjectSchema = z.object({
  actor_email: z.email({ error: "Must be a valid email address" }),
  role: z.enum([ProjectRole.member, ProjectRole.admin, ProjectRole.owner], {
    error: "Invalid role",
  }),
}) satisfies z.ZodType<AddProjectMemberRequest>;

export type MemberAddToProjectI = AddProjectMemberRequest;

export interface ProjectMemberI extends Omit<ProjectMember, "role"> {
  role: ProjectMemberRoleI;
}
