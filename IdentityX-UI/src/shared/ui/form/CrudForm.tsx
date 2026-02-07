import { useAppForm } from "@/shared/lib/forms";
import type { CrudFormConfig, FieldConfig } from "./types";

interface PropsI<TFormData> {
  formId: string;
  options: CrudFormConfig<TFormData>;
  fields: FieldConfig[];
}

export default function CrudForm<TFormData>({
  formId,
  options,
  fields
}: PropsI<TFormData>) {
    const form = useAppForm({
      defaultValues: options.defaultValues,
      validators: {
        onChange: options.validators?.onChange,
      },
      onSubmit: options.onSubmit
    });
  return (
    <form
      id={formId}
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      {fields.map(item => (
        <form.AppField
          name={item.name}
          children={(field) => (
            <field.TextField 
              label={item.label} 
              placeholder={item.placeholder}
              autoComplete={item.autoComplete}
            />
          )}
        />
      ))}
    </form>
  )
}