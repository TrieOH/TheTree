import { type AnyFormApi, type AnyFieldApi, useField } from '@tanstack/react-form';

interface FormFieldProps<TFormData, TName extends keyof TFormData & string> {
  name: TName;
  label: string;
  form: AnyFormApi; 
  children: (field: AnyFieldApi) => React.ReactNode;
}

export const FormField = <TFormData, TName extends keyof TFormData & string>({
  name,
  label,
  form,
  children,
}: FormFieldProps<TFormData, TName>) => {
  const field = useField({
    name,
    form,
  });
  const errMessage = field.state.meta.errors.map(err => err.message as string).join(", ")
  return (
    <div className="space-y-1">
      <label
        htmlFor={field.name}
        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {label}
      </label>
      {children(field)}
      {field.state.meta.errors ? (
        <p className="text-xs font-normal text-destructive">{errMessage}</p>
      ) : null}
    </div>
  );
};
