import type { 
  FieldDefinitionResultI, 
  RuleFieldCreateRequestI, 
  RuleResultI, 
  SchemaFieldCreateRequestI 
} from "../model/types";

export const mapFieldDefinitionResultToSchemaFieldCreateRequest = (fields: FieldDefinitionResultI[]) => {
  const fieldMap = new Map(fields.map(field => [field.object_id, field.key]));

  return fields.map(field => ({
    title: field.title,
    key: field.key,
    default_value: field.default_value,
    description: field.description,
    mutable: field.mutable,
    owner: field.owner,
    placeholder: field.placeholder,
    position: field.position,
    type: field.type,
    required: field.required,
    options: field.options.map(opt => ({
      label: opt.label,
      position: opt.position,
      value: opt.value
    })),
    required_rules: field.required_rules.map(rule => ({
      depends_on_field_key: fieldMap.get(rule.depends_on_field_id) ?? "",
      operator: rule.operator,
      value: rule.value,
    })),
    visibility_rules: field.visibility_rules.map(rule => ({
      depends_on_field_key: fieldMap.get(rule.depends_on_field_id) ?? "",
      operator: rule.operator,
      value: rule.value,
    }))
  })) as SchemaFieldCreateRequestI[]
}

export const mapRuleResultToRuleFieldCreateRequest = (rules: RuleResultI[], fields: FieldDefinitionResultI[]) => {
  return rules.map(rule => ({
    operator: rule.operator,
    value: rule.value,
    depends_on_field_key: fields.find(field => field.object_id === rule.depends_on_field_id)?.key
  })) as RuleFieldCreateRequestI[]
}
