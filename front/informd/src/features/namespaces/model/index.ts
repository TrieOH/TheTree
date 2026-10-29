import type {
  AddNamespaceMemberRequest,
  CreateNamespaceRequest,
  Namespace,
  NamespaceMember,
} from "@trieoh/informd-api/schemas";
import { NamespaceMemberRole } from "@trieoh/informd-api/schemas";
import z from "zod";

export const namespaceCreateSchema = z.object({
  name: z
    .string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
}) satisfies z.ZodType<CreateNamespaceRequest>;

export type NamespaceCreateI = CreateNamespaceRequest;

export type NamespaceI = Namespace;

// Member

export type NamespaceMemberRoleI =
  | typeof NamespaceMemberRole.member
  | typeof NamespaceMemberRole.admin
  | typeof NamespaceMemberRole.owner;

export const memberAddToNamespaceSchema = z.object({
  user_id: z.string({ error: "User ID is required" }),
  role: z.enum(
    [
      NamespaceMemberRole.member,
      NamespaceMemberRole.admin,
      NamespaceMemberRole.owner,
    ],
    { error: "Invalid role" },
  ),
}) satisfies z.ZodType<AddNamespaceMemberRequest>;

export type MemberAddToNamespaceI = AddNamespaceMemberRequest;

export interface NamespaceMemberI extends Omit<NamespaceMember, "role"> {
  role: NamespaceMemberRoleI;
}
