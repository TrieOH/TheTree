import type {
  CreateOrganizationRequest,
  Organization
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