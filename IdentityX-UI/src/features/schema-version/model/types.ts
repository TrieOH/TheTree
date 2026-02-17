
import z from "zod";

export const versionDraftSchema = z.object({
  schema_id: z.string(),
  project_id: z.string(),
});

export type VersionDraft = z.infer<typeof versionDraftSchema>;

export const ruleOperatorSchema = z.enum([
  "equals",
  "not_equals",
  "in",
  "not_in",
  "exists",
  "not_exists",
  "contains",
  "gt",
  "gte",
  "lt",
  "lte"
])

export type RuleOperator = z.infer<typeof ruleOperatorSchema>;

// Zod schema for Option
export const optionSchema = z.object({
  id: z.string(),
  label: z.string(),
  position: z.number(),
  value: z.string()
});

// Zod schema for Rule
export const ruleSchema = z.object({
  id: z.string(),
  depends_on_field_key: z.string().optional(),
  depends_on_field_id: z.string(),
  operator: ruleOperatorSchema,
  value: z.json()
});

export interface SchemaVersion {
  id: string;
  schema_id: string;
  based_on_version_id: string;
  status: string;
  version_number: number;
  created_at: string;
  updated_at: string;
}


const needFieldsSchema = z.object({
  default_value: z.json().optional(),
  description: z.string().optional(),
  key: z.string().min(3, "Key must be at least 3 characters long"),
  mutable: z.boolean(),
  owner: z.enum(["user", "system", "admin"]),
  position: z.number(),
  required: z.boolean(),
  placeholder: z.string().optional(),
  title: z.string().min(3, "Title must be at least 3 characters long"),
  type: z.enum(["string", "email", "int", "select", "radio", "checkbox", "bool"]),
})

export type NeedVersionField = z.infer<typeof needFieldsSchema>;

export const versionFieldSchema = needFieldsSchema.extend({
  required_rules: z.array(z.object({
    depends_on_field_key: z.string(),
    operator: ruleOperatorSchema,
    value: z.json() // is a json, but i will use string for now
  })),
  options: z.array(z.object({
    label: z.string(),
    position: z.number(),
    value: z.string()
  })),
  visibility_rules: z.array(z.object({
    depends_on_field_key: z.string(),
    operator: ruleOperatorSchema,
    value: z.json()
  }))
});

export type VersionField = z.infer<typeof versionFieldSchema>;

export type PartialVersionField = Partial<VersionField>;

export type SchemaFieldOption = z.infer<typeof versionFieldSchema.shape.options.element>;
export type SchemaFieldRequiredRule = z.infer<typeof versionFieldSchema.shape.required_rules.element>;
export type SchemaFieldVisibilityRule = z.infer<typeof versionFieldSchema.shape.visibility_rules.element>;

export const versionFieldResultZodSchema = needFieldsSchema.extend({
  id: z.string(),
  object_id: z.string().optional(),
  required_rules: z.array(ruleSchema),
  options: z.array(optionSchema),
  visibility_rules: z.array(ruleSchema)
});

export type VersionFieldResult = z.infer<typeof versionFieldResultZodSchema>;

const versionFieldListSchema = z
  .array(versionFieldResultZodSchema)
  .superRefine((fields, ctx) => {
    const existingKeys = new Set(fields.map(f => f.key));

    fields.forEach((field, fieldIndex) => {
      field.required_rules?.forEach((rule, ruleIndex) => {
        if (rule.depends_on_field_key && !existingKeys.has(rule.depends_on_field_key)) {
          ctx.addIssue({
            code: "custom",
            message: `depends_on_field_key "${rule.depends_on_field_key}" does not exist`,
            path: [
              fieldIndex,
              "required_rules",
              ruleIndex,
              "depends_on_field_key"
            ],
          });
        }
      });

      field.visibility_rules?.forEach((rule, ruleIndex) => {
        if (rule.depends_on_field_key && !existingKeys.has(rule.depends_on_field_key)) {
          ctx.addIssue({
            code: "custom",
            message: `depends_on_field_key "${rule.depends_on_field_key}" does not exist`,
            path: [
              fieldIndex,
              "visibility_rules",
              ruleIndex,
              "depends_on_field_key"
            ],
          });
        }
      });
    });
  });


export type VersionFieldList = z.infer<typeof versionFieldListSchema>;


export const schemaVersionFieldsSchema = z.object({
  schema_id: z.string(),
  project_id: z.string(),
  version: z.number(),
  fields: versionFieldListSchema,
});

export type SchemaVersionFields = z.infer<typeof schemaVersionFieldsSchema>;

export type Option = z.infer<typeof optionSchema>;

export type Rule = z.infer<typeof ruleSchema>;

export interface DetailedSchemaVersion {
  flow_id: string;
  title: string;
  fields: VersionFieldResult[];
  version_number: number;
  status: 'published' | 'draft' | 'archived';
}
