import type { FieldPath, FieldValues, UseFormReturn } from "react-hook-form";
import type { ReactNode } from "react";

/**
 * The subset of UseFormReturn that field renderers actually need.
 * Deliberately excludes `handleSubmit` / `control`, which are the only
 * members whose type depends on the zod schema's *output* (transformed)
 * type — this lets a single-generic form API be shared by every field
 * regardless of whether the schema applies `.transform()`.
 */
export type FieldFormApi<TFieldValues extends FieldValues> = Pick<
  UseFormReturn<TFieldValues>,
  "register" | "watch" | "getValues" | "setValue" | "trigger" | "resetField" | "unregister" | "formState"
>;


/**
 * Native input variants supported by the text field renderer.
 * Adding a new *string-based* variant is just a new entry here,
 * no extra component or branching logic needed.
 */
export type TextFieldInputType = "text" | "email" | "url" | "password" | "tel" | "number";

export interface TextFieldConfig<TFieldValues extends FieldValues> {
  kind: "text";
  /** Dot-path into the form values, fully typed against the zod schema. */
  name: FieldPath<TFieldValues>;
  label: string;
  placeholder?: string;
  description?: string;
  inputType?: TextFieldInputType;
  optional?: boolean;
  disabled?: boolean;
  visibleIf?: VisibilityRule<TFieldValues>;
}

export interface CustomFieldRenderArgs<TFieldValues extends FieldValues> {
  form: FieldFormApi<TFieldValues>;
}

/**
 * Escape hatch for bespoke composite UI (e.g. the social-network picker)
 * that doesn't fit a single simple input, while still living inside the
 * same step/field-list/visibility machinery as every other field.
 */
export interface CustomFieldConfig<TFieldValues extends FieldValues> {
  kind: "custom";
  /** Identifier only, used as React key. Not necessarily a form path. */
  name: string;
  visibleIf?: VisibilityRule<TFieldValues>;
  render: (args: CustomFieldRenderArgs<TFieldValues>) => ReactNode;
}

export type FieldConfig<TFieldValues extends FieldValues> =
  | TextFieldConfig<TFieldValues>
  | CustomFieldConfig<TFieldValues>;

export type FieldKind = FieldConfig<FieldValues>["kind"];

export interface StepConfig<TFieldValues extends FieldValues> {
  id: string;
  label: string;
  description?: string;
  fields: FieldConfig<TFieldValues>[];
}

/**
 * Declarative visibility rules. Keeps "show field X only if field Y has a
 * value / equals something" out of ad-hoc JSX conditionals.
 */
export type VisibilityRule<TFieldValues extends FieldValues> =
  | { type: "equals"; field: FieldPath<TFieldValues>; value: unknown }
  | { type: "notEquals"; field: FieldPath<TFieldValues>; value: unknown }
  | { type: "hasValue"; field: FieldPath<TFieldValues> }
  | { type: "predicate"; predicate: (values: TFieldValues) => boolean };