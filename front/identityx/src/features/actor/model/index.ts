import type { Actor, CreateActorRequest } from "@trieoh/identityx-api/schemas";
import { ActorAuthMethod, ActorType } from "@trieoh/identityx-api/schemas";
import z from "zod";

export type ActorAuthMethodI = "password" | "api_key";

export type ActorTypeI = "human" | "service" | "machine";

export const actorCreateSchema = z.object({
  auth_method: z.enum([ActorAuthMethod.password, ActorAuthMethod.api_key], {
    error: "Invalid auth method",
  }),
  type: z.enum([ActorType.human, ActorType.service, ActorType.machine], {
    error: "Invalid actor type",
  }),
  email: z.email({ error: "Must be a valid email address" }).optional(),
}) satisfies z.ZodType<CreateActorRequest>;

export type ActorCreateI = CreateActorRequest;

export interface ActorI extends Omit<Actor, "auth_method" | "type"> {
  auth_method: ActorAuthMethodI;
  type: ActorTypeI;
}
