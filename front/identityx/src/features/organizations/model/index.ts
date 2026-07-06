import type {
  AddOrganizationMemberRequest,
  CreateOrganizationRequest,
  Organization,
  OrganizationMember,
} from "@trieoh/identityx-models";
import {
  OrganizationRoleAdmin,
  OrganizationRoleMember,
  OrganizationRoleOwner
} from "@trieoh/identityx-models";
import z from "zod";


export const organizationCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
  slug: z.string({ error: "Slug is required" })
    .min(3, "Slug must be at least 3 characters long"),
}) satisfies z.ZodType<CreateOrganizationRequest>;

export type OrganizationCreateI = CreateOrganizationRequest;

export type OrganizationI = Organization;

// Member

export type OrganizationMemberRoleI =
  | typeof OrganizationRoleMember
  | typeof OrganizationRoleAdmin
  | typeof OrganizationRoleOwner;

export const memberAddToOrganizationSchema = z.object({
  actor_email: z.email({ error: "Must be a valid email address" }),
  role: z.enum([
    OrganizationRoleMember,
    OrganizationRoleAdmin,
    OrganizationRoleOwner
  ], { error: "Invalid role" }),
}) satisfies z.ZodType<AddOrganizationMemberRequest>;

export type MemberAddToOrganizationI = AddOrganizationMemberRequest;

export interface OrganizationMemberI
  extends Omit<OrganizationMember, "role"> {
  role: OrganizationMemberRoleI;
}
