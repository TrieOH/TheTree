import { useMemo } from "react";
import type { FieldDefinitionResultI, FieldValue } from "../../../types/fields-types";
import BasicInputField from "./BasicInputField";
import BasicSelectField from "./BasicSelectField";
import { useFieldRules } from "../../../hooks/useFieldRules";

interface DynamicFieldsProps {
  fields: FieldDefinitionResultI[];
  values: Record<string, FieldValue>;
  onValueChange: (key: string, value: FieldValue) => void;
  submitted?: boolean;
}

export default function DynamicFields({
  fields,
  values,
  onValueChange,
  submitted = false,
}: DynamicFieldsProps) {
  if (!fields || fields.length === 0) return null;

  const fieldsMap = useMemo(() => {
    return fields.reduce((acc, field) => {
      acc[field.id] = { key: field.key, title: field.title };
      if (field.object_id) acc[field.object_id] = { key: field.key, title: field.title };
      return acc;
    }, {} as Record<string, { key: string; title: string }>);
  }, [fields]);

  const sortedFields = useMemo(() => {
    return [...fields].sort((a, b) => a.position - b.position);
  }, [fields]);

  return (
    <>
      {sortedFields.map((field) => (
        <FieldRenderer
          key={field.id}
          field={field}
          values={values}
          onValueChange={onValueChange}
          submitted={submitted}
          fieldsMap={fieldsMap}
        />
      ))}
    </>
  );
}

interface FieldRendererProps {
  field: FieldDefinitionResultI;
  values: Record<string, FieldValue>;
  onValueChange: (key: string, value: FieldValue) => void;
  submitted: boolean;
  fieldsMap: Record<string, { key: string; title: string }>;
}

function FieldRenderer({
  field,
  values,
  onValueChange,
  submitted,
  fieldsMap,
}: FieldRendererProps) {
  const visibilityResult = useFieldRules(field.visibility_rules, values, fieldsMap);
  const hasVisibilityRules = field.visibility_rules && field.visibility_rules.length > 0;
  const isVisible = !hasVisibilityRules || visibilityResult.satisfied;

  const requiredResult = useFieldRules(field.required_rules, values, fieldsMap);
  const hasRequiredRules = field.required_rules && field.required_rules.length > 0;
  const isRequiredByRules = hasRequiredRules && requiredResult.satisfied;
  const isRequired = field.required || isRequiredByRules;

  const value = values[field.key] ?? field.default_value ?? "";
  const isValueEmpty = value === "" || value === undefined || value === null;

  const rulesStatus = useMemo(() => {
    const statuses = [];

    if (isRequired) statuses.push({ message: `${field.title} é obrigatório.`, passed: !isValueEmpty });

    return statuses;
  }, [
    field.title,
    isRequired,
    isValueEmpty,
  ]);

  if (!isVisible) return null;

  const commonProps = {
    name: field.key,
    label: field.title,
    placeholder: field.placeholder,
    value: String(value),
    submitted: submitted,
    rulesStatus: rulesStatus,
    required: isRequired,
  };

  switch (field.type) {
    case "string":
    case "email":
      return (
        <BasicInputField
          {...commonProps}
          type={field.type === "email" ? "email" : "text"}
          onValueChange={(val) => onValueChange(field.key, val)}
        />
      );

    case "int":
      return (
        <BasicInputField
          {...commonProps}
          type="number"
          onValueChange={(val) => onValueChange(field.key, val)}
        />
      );

    case "radio":
      return (
        <div className="trieoh trieoh-input">
          <label className="trieoh-input__label">
            {field.title}
            {isRequired && <span style={{ color: "#e74c3c", marginLeft: "4px" }}>*</span>}
          </label>
          <div className="trieoh-radio-group" style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
            {field.options?.map((opt) => (
              <label key={opt.id} style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer" }}>
                <input
                  type="radio"
                  name={field.key}
                  value={opt.value}
                  checked={value === opt.value}
                  onChange={(e) => onValueChange(field.key, e.target.value)}
                />
                <span>{opt.label}</span>
              </label>
            ))}
          </div>
          {submitted && rulesStatus
            .filter((s) => !s.passed)
            .map((s, idx) => (
              <span key={idx} style={{ color: "#e74c3c", fontSize: "12px", marginTop: "4px", display: "block" }}>
                {s.message}
              </span>
            ))}
        </div>
      );

    case "select":
      return (
        <BasicSelectField
          {...commonProps}
          options={field.options || []}
          onValueChange={(val) => onValueChange(field.key, val)}
        />
      );

    case "checkbox": {
      const currentValues = Array.isArray(values[field.key]) 
        ? (values[field.key] as string[]) 
        : [];
      
      return (
        <div className="trieoh trieoh-input">
          <label className="trieoh-input__label">
            {field.title}
            {isRequired && <span style={{ color: "#e74c3c", marginLeft: "4px" }}>*</span>}
          </label>
          <div className="trieoh-checkbox-group" style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
            {field.options?.map((opt) => (
              <label key={opt.id} style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer" }}>
                <input
                  type="checkbox"
                  name={field.key}
                  value={opt.value}
                  checked={currentValues.includes(opt.value)}
                  onChange={(e) => {
                    const newValue = e.target.checked
                      ? [...currentValues, opt.value]
                      : currentValues.filter((v) => v !== opt.value);
                    onValueChange(field.key, newValue);
                  }}
                />
                <span>{opt.label}</span>
              </label>
            ))}
          </div>
          {submitted && rulesStatus
            .filter((s) => !s.passed)
            .map((s, idx) => (
              <span key={idx} style={{ color: "#e74c3c", fontSize: "12px", marginTop: "4px", display: "block" }}>
                {s.message}
              </span>
            ))}
        </div>
      );
    }

    case "bool":
      return (
        <div 
          className="trieoh trieoh-input" 
          style={{ flexDirection: "row", alignItems: "center", gap: "8px" }}
        >
          <input
            type="checkbox"
            id={field.key}
            name={field.key}
            checked={!!values[field.key]}
            onChange={(e) => onValueChange(field.key, e.target.checked)}
          />
          <label htmlFor={field.key} className="trieoh-input__label" style={{ marginBottom: 0 }}>
            {field.title}
            {isRequired && <span style={{ color: "#e74c3c", marginLeft: "4px" }}>*</span>}
          </label>
          {submitted && rulesStatus
            .filter((s) => !s.passed)
            .map((s, idx) => (
              <span key={idx} style={{ color: "#e74c3c", fontSize: "12px", marginLeft: "8px" }}>
                {s.message}
              </span>
            ))}
        </div>
      );

    default: return null;
  }
}