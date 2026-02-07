import { useFieldContext } from "@/shared/lib/forms";
import { BasicInputField } from "@trieoh/node-auth-sdk/react";
import { useEffect, useState } from "react";

interface PropsI {
  label: string;
  placeholder: string;
  autoComplete?: string;
  submitted?: boolean;
}

export default function TextField({ label, placeholder, autoComplete, submitted }: PropsI) {
  const field = useFieldContext<string>();

  const [lastError, setLastError] = useState<string | null>(null);

  const currentErrors = field.state.meta.errors;
  const currentMessage = currentErrors.map(e => e.message).join("\n");

  useEffect(() => {
    if (currentErrors.length > 0) setLastError(currentMessage);
  }, [currentMessage, currentErrors.length]);

  const messageToShow =
    currentErrors.length > 0
      ? currentMessage
      : lastError ?? "";

  const hasError = currentErrors.length > 0;

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
          test: () => true, // not needed in this situation
        }
      ]}
    />
  )
}