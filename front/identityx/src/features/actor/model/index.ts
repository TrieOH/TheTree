import type {
  Actor,
  CreateActorRequest,
} from "@trieoh/identityx-models";
import {
  ApiKeyAuthMethod,
  HumanActorType,
  MachineActorType,
  PasswordAuthMethod,
  ServiceActorType,
} from "@trieoh/identityx-models";
import z from "zod";

export type ActorAuthMethodI =
  | "password"
  | "api_key";

export type ActorTypeI =
  | "human"
  | "service"
  | "machine";

export const actorCreateSchema = z.object({
  auth_method: z.enum([
    PasswordAuthMethod,
    ApiKeyAuthMethod,
  ], { error: "Invalid auth method" }),
  type: z.enum([
    HumanActorType,
    ServiceActorType,
    MachineActorType,
  ], { error: "Invalid actor type" }),
  email: z.email({ error: "Must be a valid email address" }).optional(),
}) satisfies z.ZodType<CreateActorRequest>;

export type ActorCreateI = CreateActorRequest;

export interface ActorI extends Omit<Actor, "auth_method" | "type"> {
  auth_method: ActorAuthMethodI;
  type: ActorTypeI;
}
