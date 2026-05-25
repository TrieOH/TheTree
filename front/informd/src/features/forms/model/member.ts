import {
  FormMemberRoleAdmin,
  FormMemberRoleEditor,
  FormMemberRoleOwner,
  FormMemberRoleViewer
} from "@trieoh/informd-models";
import type { AddFormMemberRequest, FormMember } from "@trieoh/informd-models";

import z from "zod";

export type FormMemberStatusI =
  | typeof FormMemberRoleViewer
  | typeof FormMemberRoleEditor
  | typeof FormMemberRoleAdmin
  | typeof FormMemberRoleOwner;

export const memberAddToFormSchema = z.object({
  user_id: z.string({ error: "User ID is required" }),
  role: z.enum([
    FormMemberRoleViewer,
    FormMemberRoleEditor,
    FormMemberRoleAdmin,
    FormMemberRoleOwner
  ], { error: "Invalid role" }),
}) satisfies z.ZodType<AddFormMemberRequest>;

export type MemberAddToFormI = AddFormMemberRequest;

export interface FormMemberI
  extends Omit<FormMember, "role"> {
  role: FormMemberStatusI;
}