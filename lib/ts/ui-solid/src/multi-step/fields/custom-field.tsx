import type { JSX } from "@solidjs/web";
import { Field } from "../../input";
import type { CustomMultiStepField } from "../types";

export interface MultiStepCustomFieldProps<T = Record<string, unknown>> {
  field: CustomMultiStepField<T>;
  value?: unknown;
  error?: string;
  onChange?: (value: unknown) => void;
}

export function MultiStepCustomField<T = Record<string, unknown>>(
  props: MultiStepCustomFieldProps<T>,
): JSX.Element {
  return (
    <Field
      label={props.field.label}
      required={props.field.required}
      hint={props.field.hint}
      error={props.error}
    >
      {(ids) =>
        props.field.render({
          value: props.value,
          onChange: (v) => props.onChange?.(v),
          error: props.error,
          ids,
        })
      }
    </Field>
  );
}
