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


const needFieldsSchema = z.object({
  default_value: z.json().optional(), // done
  description: z.string().optional(), // done
  key: z.string().min(3, "Key must be at least 3 characters long"), // done
  mutable: z.boolean(), // done
  owner: z.enum(["user", "system", "admin"]), // done
  position: z.number(), // done
  required: z.boolean(), // done
  placeholder: z.string().optional(), // done
  title: z.string().min(3, "Title must be at least 3 characters long"), // done
  type: z.enum(["string", "email", "int", "select", "radio", "checkbox", "bool"]), // done
})

type NeedVersionField = z.infer<typeof needFieldsSchema>;

export const versionFieldSchema = needFieldsSchema.extend({
  required_rules: z.array(z.object({
    depends_on_field_key: z.string(),
    operator: ruleOperatorSchema,
    value: z.json() // is a json, but i will use string for now
  })),
  options: z.array(z.object({ // done
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


export type VersionFieldResult = NeedVersionField & {
  object_id: string;
  required_rules: {
    id: string;
    depends_on_field_id: string;
    operator: RuleOperator;
    value: string;
  }[],
  options: {
    id: string;
    label: string,
    position: number,
    value: string
  }[],
  visibility_rules: {
    id: string;
    depends_on_field_id: string;
    operator: RuleOperator;
    value: string;
  }[]; 
}

export interface DetailedSchemaVersion {
  flow_id: string;
  title: string;
  fields: VersionFieldResult[];
  version_number: number;
}

