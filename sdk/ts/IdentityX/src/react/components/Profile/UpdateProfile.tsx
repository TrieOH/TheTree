import { useState, useEffect, type MouseEvent } from "react";
import { useAuth } from "../../AuthProvider";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import type { FieldDefinitionResultI, FieldValue } from "../../../types/fields-types";
import DynamicFields from "../Form/DynamicFields";

export interface UpdateProfileProps {
  fields: FieldDefinitionResultI[];
  initialValues?: Record<string, FieldValue>;
  onSuccess?: () => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  submitLabel?: string;
}

export function UpdateProfile({
  fields,
  initialValues,
  onSuccess,
  onFailed,
  submitLabel = "Atualizar Dados",
}: UpdateProfileProps) {
  const [values, setValues] = useState<Record<string, FieldValue>>(initialValues || {});
  const [submitted, setSubmitted] = useState(false);
  const [loading, setLoading] = useState(false);
  const { auth } = useAuth();

  useEffect(() => {
    if (initialValues) {
      setValues(initialValues);
    }
  }, [initialValues]);

  const handleValueChange = (key: string, value: FieldValue) => {
    setValues(prev => ({ ...prev, [key]: value }));
  };

  const handleSubmit = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setSubmitted(true);

    const invalid = fields.some(f => f.required && !values[f.key]);
    if (invalid) return;

    setLoading(true);
    const res = await auth.updateProfile(values);
    if (res.code === 201 && onSuccess) await onSuccess();
    else if (onFailed) await onFailed(res.message, res.trace);
    setLoading(false);
  };

  return (
    <form className="trieoh trieoh-card trieoh-card--full-rounded">
      <div className="trieoh-card__fields">
        <DynamicFields 
          fields={fields}
          values={values}
          onValueChange={handleValueChange}
          submitted={submitted}
        />
      </div>
      <BasicSubmitButton 
        label={submitLabel} 
        onSubmit={handleSubmit} 
        loading={loading}
      />
    </form>
  );
}
