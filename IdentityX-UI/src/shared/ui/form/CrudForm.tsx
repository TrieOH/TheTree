import { useAppForm } from "@/shared/lib/forms";
import type { CrudFormConfig, FieldConfig } from "./types";
import { useState } from "react";

interface PropsI<TFormData> {
  formId: string;
  options: CrudFormConfig<TFormData>;
  fields: FieldConfig[];
}

export default function CrudForm<TFormData>({
  formId,
  options,
  fields,
}: PropsI<TFormData>) {
  const [submitted, setSubmitted] = useState(false);
  const form = useAppForm({
    defaultValues: options.defaultValues,
    validators: options.validators,
    onSubmit: async ({ value, formApi }) => {
      formApi.reset()
      formApi.mount()
      await options.onSubmit({ value, formApi })
      setSubmitted(false);
    }
  });
  return (
    <form
      id={formId}
      onSubmit={async (e) => {
        e.preventDefault();
        e.stopPropagation();
        setSubmitted(true);
        await form.handleSubmit();
      }}
    >
      {fields.map(item => (
        <form.AppField
          key={item.name}
          name={item.name}
        >
          {(field) => (
            <field.TextField 
              label={item.label} 
              placeholder={item.placeholder}
              autoComplete={item.autoComplete}
              errors={item.errors}
              submitted={submitted}
            />
          )}
        </form.AppField>
      ))}
    </form>
  )
}