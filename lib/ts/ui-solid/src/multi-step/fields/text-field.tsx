import type { JSX } from "@solidjs/web";
import { Field, Input } from "../../input";
import type { TextMultiStepField } from "../types";

export interface MultiStepTextFieldProps<T = Record<string, unknown>> {
  field: TextMultiStepField<T>;
  value?: unknown;
  error?: string;
  onChange?: (value: string) => void;
}

export function MultiStepTextField<T = Record<string, unknown>>(
  props: MultiStepTextFieldProps<T>,
): JSX.Element {
  const inputType = () => props.field.kind ?? "text";
  const strValue = () => (props.value != null ? String(props.value) : "");

  return (
    <Field
      label={props.field.label}
      required={props.field.required}
      hint={props.field.hint}
      error={props.error}
    >
      {(ids) => (
        <Input
          id={ids.id}
          type={inputType()}
          aria-describedby={ids.describedBy}
          aria-invalid={props.error ? "true" : undefined}
          value={strValue()}
          placeholder={props.field.placeholder}
          disabled={props.field.disabled}
          autocomplete="off"
          onInput={(event) => props.onChange?.(event.currentTarget.value)}
        />
      )}
    </Field>
  );
}
