import type { VersionField, VersionFieldResult } from "../model/types";

export const mapFieldIdsToKeys = (fields: VersionFieldResult[]): VersionFieldResult[] => {
  const fieldMap = new Map(
    fields.map(field => [field.object_id, field.key])
  );

  return fields.map(field => ({
    ...field,
    required_rules: field.required_rules.map(rule => ({
      id: rule.id,
      operator: rule.operator,
      value: rule.value,
      depends_on_field_id: rule.depends_on_field_id,
      depends_on_field_key: fieldMap.get(rule.depends_on_field_id) ?? ""
    })),
    visibility_rules: field.visibility_rules.map(rule => ({
      id: rule.id,
      operator: rule.operator,
      value: rule.value,
      depends_on_field_id: rule.depends_on_field_id,
      depends_on_field_key: fieldMap.get(rule.depends_on_field_id) ?? ""
    }))
  }));
};

export const mapVersionFieldResultToVersionField = (field: VersionFieldResult): VersionField => {
  return {
    key: field.key,
    mutable: field.mutable,
    options: field.options,
    owner: field.owner,
    position: field.position,
    required: field.required,
    required_rules: field.required_rules.map(rule => ({
      depends_on_field_key: rule.depends_on_field_key || '',
      operator: rule.operator,
      value: rule.value
    })),
    visibility_rules: field.visibility_rules.map(rule => ({
      depends_on_field_key: rule.depends_on_field_key || '',
      operator: rule.operator,
      value: rule.value
    })),
    default_value: field.default_value,
    description: field.description,
    placeholder: field.placeholder,
    title: field.title,
    type: field.type,
  }
};