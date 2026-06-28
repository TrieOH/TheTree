import { useFieldContext } from "@/shared/lib/forms";
import { BasicInputField } from "@trieoh/identityx-sdk-ts/react";
import type { RuleStatus } from "./types";

interface PropsI {
  label: string;
  placeholder: string;
  autoComplete?: string;
  required?: boolean;
  getRulesStatus?: (value: unknown) => RuleStatus[]
  submitted?: boolean;
}

export default function TextField({ label, placeholder, autoComplete, required, getRulesStatus, submitted }: PropsI) {
  const field = useFieldContext<string>();
  return (
    <BasicInputField
      name={field.name}
      label={required ? `${label} *` : label}
      placeholder={placeholder}
      value={field.state.value}
      onValueChange={(v) => field.handleChange(v)}
      onBlur={field.handleBlur}
      autoComplete={autoComplete}
      submitted={submitted}
      rulesStatus={getRulesStatus ? getRulesStatus(field.state.value) : []}
    />
  )
}