import { EventMemberRole } from "@trieoh/univents-api/schemas";
import z from "zod";

export const eventMemberRoles = [
  EventMemberRole.owner,
  EventMemberRole.admin,
  EventMemberRole.staff,
] as const;

export const eventMemberCreateSchema = z.object({
  email: z.email("Informe um e-mail válido"),
  role: z.enum(eventMemberRoles),
});

export type EventMemberCreateInput = z.input<typeof eventMemberCreateSchema>;
export type EventMemberCreateOutput = z.output<typeof eventMemberCreateSchema>;

export { EventMemberRole };
