import type { Capability, CreateCapabilityRequest } from "@trieoh/identityx-models";
import z from "zod";

export const capabilityCreateSchema = z.object({
  action: z.string().min(3, "Action must be at least 3 characters long"),
  resource: z.string().min(3, "Resource must be at least 3 characters long"),
}) satisfies z.ZodType<CreateCapabilityRequest>;

export type CapabilityCreateI = z.infer<typeof capabilityCreateSchema>;

export type CapabilityI = Capability;
