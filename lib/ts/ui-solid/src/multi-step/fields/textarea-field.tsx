import type { JSX } from "@solidjs/web";
import { Field, Textarea } from "../../input";
import type { TextareaMultiStepField } from "../types";

export interface MultiStepTextareaFieldProps<T = Record<string, unknown>> {
  field: TextareaMultiStepField<T>;
  value?: unknown;
  error?: string;
  onChange?: (value: string) => void;
}

export function MultiStepTextareaField<T = Record<string, unknown>>(
  props: MultiStepTextareaFieldProps<T>,
): JSX.Element {
  const strValue = () => (props.value != null ? String(props.value) : "");

  return (
    <Field
      label={props.field.label}
      required={props.field.required}
      hint={props.field.hint}
      error={props.error}
    >
      {(ids) => (
        <Textarea
          id={ids.id}
          rows={props.field.rows ?? 4}
          aria-describedby={ids.describedBy}
          aria-invalid={props.error ? "true" : undefined}
          value={strValue()}
          placeholder={props.field.placeholder}
          disabled={props.field.disabled}
          onInput={(event) => props.onChange?.(event.currentTarget.value)}
        />
      )}
    </Field>
  );
}
