import {
  NamespaceMemberRoleAdmin,
  NamespaceMemberRoleEditor,
  NamespaceMemberRoleOwner,
  NamespaceMemberRoleViewer
} from "@trieoh/informd-models";
import type {
  AddNamespaceMemberRequest,
  CreateFormRequest,
  CreateNamespaceRequest,
  Form,
  FormStatusArchived,
  FormStatusClosed,
  FormStatusDraft,
  FormStatusOpen,
  Namespace,
  NamespaceMember
} from "@trieoh/informd-models";
import z from "zod";

export const namespaceCreateSchema = z.object({
  name: z.string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
}) satisfies z.ZodType<CreateNamespaceRequest>;

export type NamespaceCreateI = CreateNamespaceRequest;

export type NamespaceI = Namespace;

// Member

export type NamespaceMemberRoleI =
  | typeof NamespaceMemberRoleViewer
  | typeof NamespaceMemberRoleEditor
  | typeof NamespaceMemberRoleAdmin
  | typeof NamespaceMemberRoleOwner;

export const memberAddToNamespaceSchema = z.object({
  user_id: z.string({ error: "User ID is required" }),
  role: z.enum([
    NamespaceMemberRoleViewer,
    NamespaceMemberRoleEditor,
    NamespaceMemberRoleAdmin,
    NamespaceMemberRoleOwner
  ], { error: "Invalid role" }),
}) satisfies z.ZodType<AddNamespaceMemberRequest>;

export type MemberAddToNamespaceI = AddNamespaceMemberRequest;

export interface NamespaceMemberI
  extends Omit<NamespaceMember, "role"> {
  role: NamespaceMemberRoleI;
}

// Form

export type FormStatusI =
  | typeof FormStatusDraft
  | typeof FormStatusOpen
  | typeof FormStatusClosed
  | typeof FormStatusArchived;

export const formCreateOnNamespaceSchema = z.object({
  title: z.string({ error: "Title is required" })
    .min(3, "Title must be at least 3 characters long"),
}) satisfies z.ZodType<CreateFormRequest>;

export type FormCreateOnNamespaceI = CreateFormRequest;

export type FormI = Form;
