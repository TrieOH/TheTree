import { useFieldContext } from "@/shared/lib/forms";
import { BasicSelectField } from "@trieoh/node-auth-sdk/react";


interface Option {
  label: string;
  value: string;
}

interface PropsI {
  label: string;
  placeholder?: string;
  options: Option[];
  errors?: string[];
  submitted?: boolean;
}

export default function SelectField({ label, placeholder, options, errors, submitted }: PropsI) {
  const field = useFieldContext<string>();

  const currentErrors = field.state.meta.errors;
  const hasError = currentErrors.length > 0;
  const messageToShow = errors?.join("\n") || "";


  return (
    <BasicSelectField
      label={label}
      name={field.name}
      placeholder={placeholder}
      value={field.state.value}
      onValueChange={(v) => field.handleChange(v)}
      onBlur={field.handleBlur}
      submitted={submitted}
      options={options.map(opt => ({
        id: opt.value,
        label: opt.label,
        value: opt.value
      }))}
      rulesStatus={[
        {
          message: messageToShow,
          passed: !hasError,
        }
      ]}
    />
  );
}
