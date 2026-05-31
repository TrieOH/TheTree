import {
  FieldTypeBool,
  FieldTypeDate,
  FieldTypeDatetime,
  FieldTypeEmail,
  FieldTypeFile,
  FieldTypeFloat,
  FieldTypeInt,
  FieldTypePhone,
  FieldTypeSelect,
  FieldTypeString,
  FieldTypeTime,
  FieldTypeURL,
  SelectBehaviourCheckbox,
  SelectBehaviourDropdownCheckbox,
  SelectBehaviourDropdownRadio,
  SelectBehaviourRadio,
  SelectValueTypeDate,
  SelectValueTypeDatetime,
  SelectValueTypeEmail,
  SelectValueTypeFloat,
  SelectValueTypeInt,
  SelectValueTypePhone,
  SelectValueTypeString,
  SelectValueTypeTime,
  SelectValueTypeURL,
} from "@trieoh/informd-models";
import type { CreateFieldRequest, CreateFieldSelectConfigRequest, Field, FieldSelectConfig, UpdateFieldRequest } from "@trieoh/informd-models";

import z from "zod";

export type FieldTypeI =
  | typeof FieldTypeString
  | typeof FieldTypeEmail
  | typeof FieldTypeInt
  | typeof FieldTypeFloat
  | typeof FieldTypeBool
  | typeof FieldTypeDate
  | typeof FieldTypeTime
  | typeof FieldTypeDatetime
  | typeof FieldTypeSelect
  | typeof FieldTypeFile
  | typeof FieldTypePhone
  | typeof FieldTypeURL;

export type SelectBehaviourI =
  | typeof SelectBehaviourCheckbox
  | typeof SelectBehaviourRadio
  | typeof SelectBehaviourDropdownCheckbox
  | typeof SelectBehaviourDropdownRadio;

export type SelectValueTypeI =
  | typeof SelectValueTypeString
  | typeof SelectValueTypeEmail
  | typeof SelectValueTypeInt
  | typeof SelectValueTypeFloat
  | typeof SelectValueTypeDate
  | typeof SelectValueTypeTime
  | typeof SelectValueTypeDatetime
  | typeof SelectValueTypePhone
  | typeof SelectValueTypeURL;

/** Schema for creating/updating a select config */
const createFieldSelectConfigSchema = z.object({
  behaviour: z.enum([
    SelectBehaviourCheckbox,
    SelectBehaviourRadio,
    SelectBehaviourDropdownCheckbox,
    SelectBehaviourDropdownRadio
  ], { error: "Invalid select behaviour" }),
  value_type: z.enum([
    SelectValueTypeString,
    SelectValueTypeEmail,
    SelectValueTypeInt,
    SelectValueTypeFloat,
    SelectValueTypeDate,
    SelectValueTypeTime,
    SelectValueTypeDatetime,
    SelectValueTypePhone,
    SelectValueTypeURL
  ], { error: "Invalid select value type" }),
  options: z.any(),
}) satisfies z.ZodType<CreateFieldSelectConfigRequest>;


export const createFieldRequestSchema = z.object({
  key: z.string({ error: "Field key is required" }),
  title: z.string({ error: "Field title is required" }),
  description: z.string().optional(),
  position_hint: z.number({ error: "Position hint must be a number" }),
  required: z.boolean({ error: "Required must be a boolean" }),
  type: z.enum([
    FieldTypeString,
    FieldTypeEmail,
    FieldTypeInt,
    FieldTypeFloat,
    FieldTypeBool,
    FieldTypeDate,
    FieldTypeTime,
    FieldTypeDatetime,
    FieldTypeSelect,
    FieldTypeFile,
    FieldTypePhone,
    FieldTypeURL
  ], { error: "Invalid type" }),
  placeholder: z.any().optional(),
  default_value: z.any().optional(),
  config: z.any().optional(),
  select_config: createFieldSelectConfigSchema.optional(), // Only required when type is "select"
}) satisfies z.ZodType<CreateFieldRequest>;

export const fieldUpdateRequestSchema = z.object({
  id: z.string({ error: "Field ID is required" }),
  key: z.string({ error: "Field key is required" }),
  title: z.string({ error: "Field title is required" }),
  description: z.string().optional(),
  position_hint: z.number({ error: "Position hint must be a number" }),
  required: z.boolean({ error: "Required must be a boolean" }),
  type: z.enum([
    FieldTypeString,
    FieldTypeEmail,
    FieldTypeInt,
    FieldTypeFloat,
    FieldTypeBool,
    FieldTypeDate,
    FieldTypeTime,
    FieldTypeDatetime,
    FieldTypeSelect,
    FieldTypeFile,
    FieldTypePhone,
    FieldTypeURL
  ], { error: "Invalid type" }),
  placeholder: z.any().optional(),
  default_value: z.any().optional(),
  config: z.any().optional(),
  select_config: createFieldSelectConfigSchema.optional(),
}) satisfies z.ZodType<UpdateFieldRequest>;

export type CreateFieldRequestI = CreateFieldRequest;
export type FieldUpdateI = UpdateFieldRequest;
export type CreateFieldSelectConfigRequestI = CreateFieldSelectConfigRequest;

export interface FieldI extends Omit<Field, "type"> {
  type: FieldTypeI;
}

export interface FieldSelectConfigI extends Omit<FieldSelectConfig, "behaviour" | "value_type"> {
  behaviour: SelectBehaviourI;
  value_type: SelectValueTypeI;
}