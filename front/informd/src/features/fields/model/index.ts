import type {
  CreateFieldRequest,
  CreateFieldSelectConfigRequest,
  Field,
  FieldSelectConfig,
  UpdateFieldRequest,
} from "@trieoh/informd-api/schemas";
import {
  FieldType,
  SelectBehaviour,
  SelectValueType,
} from "@trieoh/informd-api/schemas";

import z from "zod";
import type { RawJSON } from "#/shared/lib/raw-json";

export type FieldTypeI = FieldType;

export type SelectBehaviourI = SelectBehaviour;

export type SelectValueTypeI = SelectValueType;

/** Schema for creating/updating a select config */
const createFieldSelectConfigSchema = z.object({
  behaviour: z.enum(
    [
      SelectBehaviour.checkbox,
      SelectBehaviour.radio,
      SelectBehaviour["dropdown-checkbox"],
      SelectBehaviour["dropdown-radio"],
    ],
    { error: "Invalid select behaviour" },
  ),
  value_type: z.enum(
    [
      SelectValueType.string,
      SelectValueType.email,
      SelectValueType.int,
      SelectValueType.float,
      SelectValueType.date,
      SelectValueType.time,
      SelectValueType.datetime,
      SelectValueType.phone,
      SelectValueType.url,
    ],
    { error: "Invalid select value type" },
  ),
  options: z.any(),
}) satisfies z.ZodType<CreateFieldSelectConfigRequestI>;

export const createFieldRequestSchema = z.object({
  key: z.string({ error: "Field key is required" }),
  title: z.string({ error: "Field title is required" }),
  description: z.string().optional(),
  position_hint: z.number({ error: "Position hint must be a number" }),
  required: z.boolean({ error: "Required must be a boolean" }),
  type: z.enum(
    [
      FieldType.string,
      FieldType.email,
      FieldType.int,
      FieldType.float,
      FieldType.bool,
      FieldType.date,
      FieldType.time,
      FieldType.datetime,
      FieldType.select,
      FieldType.file,
      FieldType.phone,
      FieldType.url,
    ],
    { error: "Invalid type" },
  ),
  placeholder: z.any().optional(),
  default_value: z.any().optional(),
  config: z.any().optional(),
  select_config: createFieldSelectConfigSchema.optional(), // Only required when type is "select"
}) satisfies z.ZodType<CreateFieldRequestI>;

export const fieldUpdateRequestSchema = z.object({
  id: z.string({ error: "Field ID is required" }),
  key: z.string({ error: "Field key is required" }),
  title: z.string({ error: "Field title is required" }),
  description: z.string().optional(),
  position_hint: z.number({ error: "Position hint must be a number" }),
  required: z.boolean({ error: "Required must be a boolean" }),
  type: z.enum(
    [
      FieldType.string,
      FieldType.email,
      FieldType.int,
      FieldType.float,
      FieldType.bool,
      FieldType.date,
      FieldType.time,
      FieldType.datetime,
      FieldType.select,
      FieldType.file,
      FieldType.phone,
      FieldType.url,
    ],
    { error: "Invalid type" },
  ),
  placeholder: z.any().optional(),
  default_value: z.any().optional(),
  config: z.any().optional(),
  select_config: createFieldSelectConfigSchema.optional(),
}) satisfies z.ZodType<FieldUpdateI>;

type RawFields = "placeholder" | "default_value" | "config";

export interface CreateFieldSelectConfigRequestI
  extends Omit<CreateFieldSelectConfigRequest, "options"> {
  options: RawJSON;
}

export interface CreateFieldRequestI
  extends Omit<CreateFieldRequest, RawFields | "select_config"> {
  placeholder?: RawJSON;
  default_value?: RawJSON;
  config?: RawJSON;
  select_config?: CreateFieldSelectConfigRequestI;
}

export interface FieldUpdateI
  extends Omit<UpdateFieldRequest, RawFields | "select_config"> {
  placeholder?: RawJSON;
  default_value?: RawJSON;
  config?: RawJSON;
  select_config?: CreateFieldSelectConfigRequestI;
}

export interface FieldI extends Omit<Field, "type" | RawFields> {
  type: FieldTypeI;
  placeholder?: RawJSON;
  default_value?: RawJSON;
  config?: RawJSON;
}

export interface OptionI {
  label: string;
  value: string;
}

export interface FieldSelectConfigI
  extends Omit<FieldSelectConfig, "behaviour" | "value_type" | "options"> {
  behaviour: SelectBehaviourI;
  value_type: SelectValueTypeI;
  options: (OptionI | string)[];
}
