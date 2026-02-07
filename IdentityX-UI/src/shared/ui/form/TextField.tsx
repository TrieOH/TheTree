import { useFieldContext } from "@/shared/lib/forms";
import { BasicInputField } from "@trieoh/node-auth-sdk/react";

interface PropsI {
  label: string;
  placeholder: string;
  autoComplete?: string;
}

export default function TextField({ label, placeholder, autoComplete }: PropsI) {
  const field = useFieldContext<string>();
  return (
    <BasicInputField 
      name={field.name}
      label={label}
      placeholder={placeholder}
      value={field.state.value}
      onValueChange={(v) => field.handleChange(v)}
      onBlur={field.handleBlur}
      autoComplete={autoComplete}
    />
  )
}