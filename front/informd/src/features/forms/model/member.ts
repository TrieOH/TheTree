import type {
  AddFormMemberRequest,
  FormMember,
} from "@trieoh/informd-api/schemas";
import { FormMemberRole } from "@trieoh/informd-api/schemas";

import z from "zod";

export type FormMemberStatusI =
  | typeof FormMemberRole.member
  | typeof FormMemberRole.admin
  | typeof FormMemberRole.owner;

export const memberAddToFormSchema = z.object({
  user_id: z.string({ error: "User ID is required" }),
  role: z.enum(
    [FormMemberRole.member, FormMemberRole.admin, FormMemberRole.owner],
    { error: "Invalid role" },
  ),
}) satisfies z.ZodType<AddFormMemberRequest>;

export type MemberAddToFormI = AddFormMemberRequest;

export interface FormMemberI extends Omit<FormMember, "role"> {
  role: FormMemberStatusI;
}
