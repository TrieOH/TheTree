import { useState } from "react";
import type { FieldDefinitionResultI, FieldValue } from "../../../types/fields-types";
import DynamicFields from "./DynamicFields";

export interface EditorFormProps {
  fields: FieldDefinitionResultI[];
  submitted?: boolean;
  title?: string;
  description?: string;
  noFieldsMessage?: string;
}

export function EditorForm({
  fields,
  submitted = false,
  title,
  description,
  noFieldsMessage = "No fields available",
}: EditorFormProps) {
  const [values, setValues] = useState<Record<string, FieldValue>>({});

  const handleValueChange = (key: string, value: FieldValue) => {
    setValues((prev) => ({ ...prev, [key]: value }));
  };

  return (
    <div 
      className="trieoh trieoh-card trieoh-card--full-rounded"
      style={{ gap: '0.25rem' }}
    >
      {(title || description) && (
        <div style={{ width: "100%", textAlign: 'center'}}>
          {title && <h3 style={{ margin: 0, fontSize: "1.25rem", fontWeight: 600 }}>{title}</h3>}
          {description && <p style={{ margin: "0.5rem 0 0", fontSize: "0.875rem", opacity: 0.7 }}>{description}</p>}
          <hr style={{ marginTop: "1rem", border: 0, borderTop: "1px solid rgba(0,0,0,0.1)", marginBottom: "0.5rem" }} />
        </div>
      )}
      
      <div className="trieoh-card__fields" >
        {fields && fields.length > 0 ? (
          <DynamicFields 
            fields={fields}
            values={values}
            onValueChange={handleValueChange}
            submitted={submitted}
          />
        ) : (
          <div className="trieoh-card__empty">
            {noFieldsMessage}
          </div>
        )}
      </div>
    </div>
  );
}
