import { useFieldContext } from "@/shared/lib/forms";
import { BasicInputField } from "@trieoh/identityx-sdk-ts/react";

interface PropsI {
  label: string;
  placeholder: string;
  autoComplete?: string;
  errors?: string[];
  submitted?: boolean;
}

export default function TextField({ label, placeholder, autoComplete, errors, submitted }: PropsI) {
  const field = useFieldContext<string>();

  const currentErrors = field.state.meta.errors;
  const hasError = currentErrors.length > 0;
  const messageToShow = errors?.join("\n") || "";

  return (
    <BasicInputField 
      name={field.name}
      label={label}
      placeholder={placeholder}
      value={field.state.value}
      onValueChange={(v) => field.handleChange(v)}
      onBlur={field.handleBlur}
      autoComplete={autoComplete}
      submitted={submitted}
      rulesStatus={[
        {
          message: messageToShow,
          passed: !hasError,
        }
      ]}
    />
  )
}