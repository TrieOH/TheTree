export type DevStatus = "draft" | "published" | "archived";

export const OPERATORS = [
  "equals",
  "not_equals",
  "in",
  "not_in",
  "exists",
  "not_exists",
  "greater_than",
  "greater_than_equal",
  "lower_than",
  "lower_than_equal",
  "contains",
] as const;

export type Operator = typeof OPERATORS[number];

export type FieldValue = string | number | boolean | string[] | undefined;

export interface RuleResultI {
  depends_on_field_id: string;
  id: string;
  operator: string;
  value: FieldValue;
}

export interface OptionResultI {
  id: string;
  label: string;
  value: string;
  position: number;
}

export interface FieldDefinitionResultI {
  id: string;
  object_id: string;
  key: string;
  title: string;
  type: "string" | "email" | "int" | "select" | "radio" | "checkbox" | "bool";
  placeholder: string;
  description: string;
  position: number;
  options: OptionResultI[];
  default_value: FieldValue;
  mutable: boolean;
  required: boolean;
  owner: "user" | "admin" | "system"
  visibility_rules: RuleResultI[];
  required_rules: RuleResultI[];
  created_at: string;
  updated_at: string;
}

export interface ProjectFieldDefinitionResultI {
  id: string;
  flow_id: string;
  schema_id: string;
  version_id: string;
  title: string;
  schema_type: "context" | "sub-context" | "core";
  status: DevStatus;
  fields: FieldDefinitionResultI[];
  version_number: number;
  created_at: string;
  updated_at: string;
}