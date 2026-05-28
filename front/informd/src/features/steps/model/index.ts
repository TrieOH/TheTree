import type {
  CreateStepRequest,
  Step,
  UpdateStepRequest
} from "@trieoh/informd-models";
import z from "zod";

export const stepCreateSchema = z.object({
  title: z.string({ error: "Title is required" })
    .min(3, "Title must be at least 3 characters long"),
  description: z.string().optional(),
  position_hint: z.number({ error: "Position hint is required" })
}) satisfies z.ZodType<CreateStepRequest>;

export type StepCreateI = CreateStepRequest;

export const stepUpdateBulkSchema = z.object({
  title: z.string({ error: "Title is required" })
    .min(3, "Title must be at least 3 characters long"),
  description: z.string().optional(),
  position_hint: z.number({ error: "Position hint is required" })
}) satisfies z.ZodType<CreateStepRequest>;

export type StepUpdateI = UpdateStepRequest;

export type StepI = Step;