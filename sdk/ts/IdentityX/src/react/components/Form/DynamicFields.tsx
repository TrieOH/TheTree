import type { FieldDefinitionResultI, FieldValue } from "../../../types/fields-types";
import BasicInputField from "./BasicInputField";

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

  return (
    <>
      {fields
        .sort((a, b) => a.position - b.position)
        .map((field) => {
          const value = values[field.key] ?? field.default_value ?? "";
          const isValueEmpty = value === "" || value === undefined || value === null;
          
          const commonProps = {
            key: field.id,
            name: field.key,
            label: field.title,
            placeholder: field.placeholder,
            value: String(value),
            submitted: submitted,
            rulesStatus: field.required && isValueEmpty ? [
                { message: `${field.title} é obrigatório.`, passed: false }
            ] : field.required ? [
                { message: `${field.title} é obrigatório.`, passed: true }
            ] : [],
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
                <div key={field.id} className="trieoh trieoh-input">
                  <label className="trieoh-input__label">{field.title}</label>
                  <div className="trieoh-radio-group" style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    {field.options.map((opt) => (
                      <label key={opt.id} style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
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
                </div>
              );
            case "select":
              return (
                <div key={field.id} className="trieoh trieoh-input">
                  <label htmlFor={field.key} className="trieoh-input__label">
                    {field.title}
                  </label>
                  <div className="trieoh-input__container">
                    <select
                      id={field.key}
                      name={field.key}
                      value={String(values[field.key] ?? field.default_value ?? "")}
                      onChange={(e) => onValueChange(field.key, e.target.value)}
                      className="trieoh-input__container-field"
                      style={{ width: '100%', background: 'transparent', border: 'none', outline: 'none', color: 'inherit' }}
                    >
                      <option value="" disabled>{field.placeholder || "Selecione uma opção"}</option>
                      {field.options.map((opt) => (
                        <option key={opt.id} value={opt.value}>
                          {opt.label}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>
              );
            case "checkbox":
              const currentValues = Array.isArray(values[field.key]) ? (values[field.key] as string[]) : [];
              return (
                <div key={field.id} className="trieoh trieoh-input">
                  <label className="trieoh-input__label">{field.title}</label>
                  <div className="trieoh-checkbox-group" style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    {field.options.map((opt) => (
                      <label key={opt.id} style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
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
                </div>
              );
            case "bool":
              return (
                <div key={field.id} className="trieoh trieoh-input" style={{ flexDirection: 'row', alignItems: 'center', gap: '8px' }}>
                  <input
                    type="checkbox"
                    id={field.key}
                    name={field.key}
                    checked={!!values[field.key]}
                    onChange={(e) => onValueChange(field.key, e.target.checked)}
                  />
                  <label htmlFor={field.key} className="trieoh-input__label" style={{ marginBottom: 0 }}>
                    {field.title}
                  </label>
                </div>
              );
            default:
              return null;
          }
        })}
    </>
  );
}
