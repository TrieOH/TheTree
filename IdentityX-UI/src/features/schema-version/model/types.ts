import z from "zod";

export const versionDraftSchema = z.object({
  schema_id: z.string(),
  project_id: z.string(),
});

export type VersionDraft = z.infer<typeof versionDraftSchema>;

export const ruleOperatorSchema = z.enum(["equals", "not_equals", "in", "not_in", "exists", "not_exists"])

export type RuleOperator = z.infer<typeof ruleOperatorSchema>;


export interface SchemaVersion {
  id: string;
  schema_id: string;
  based_on_version_id: string;
  status: string;
  version_number: number;
  created_at: string;
  updated_at: string;
}

export const versionFieldSchema = z.object({
  default_value: z.json().optional(), // done
  description: z.string().optional(), // done
  key: z.string().min(3, "Key must be at least 3 characters long"), // done
  mutable: z.boolean(), // done
  options: z.array(z.object({ // done
    label: z.string(),
    position: z.number(),
    value: z.string()
  })),
  owner: z.enum(["user", "system", "admin"]), // done
  placeholder: z.string().optional(), // done
  position: z.number(), // done
  required: z.boolean(), // done
  required_rules: z.array(z.object({
    depends_on_field_key: z.string(),
    operator: ruleOperatorSchema,
    value: z.json() // is a json, but i will use string for now
  })),
  title: z.string().min(3, "Title must be at least 3 characters long"), // done
  type: z.enum(["string", "email", "int", "select", "radio", "checkbox", "bool"]), // done
  visibility_rules: z.array(z.object({
    depends_on_field_key: z.string(),
    operator: ruleOperatorSchema,
    value: z.json()
  }))
});

export type VersionField = z.infer<typeof versionFieldSchema>;

const versionFieldListSchema = z
  .array(versionFieldSchema)
  .superRefine((fields, ctx) => {
    const existingKeys = new Set(fields.map(f => f.key));

    fields.forEach((field, fieldIndex) => {
      field.required_rules?.forEach((rule, ruleIndex) => {
        if (!existingKeys.has(rule.depends_on_field_key)) {
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
        if (!existingKeys.has(rule.depends_on_field_key)) {
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