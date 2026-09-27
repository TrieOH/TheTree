import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { Label, Textarea } from "../../input";
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
  const id = `field-${props.field.name}`;
  const errorId = `${id}-error`;
  const hintId = `${id}-hint`;

  const describedBy = () => {
    if (props.error) return errorId;
    if (props.field.hint) return hintId;
    return undefined;
  };

  const handleInput = (event: InputEvent & { currentTarget: HTMLTextAreaElement }) => {
    const val = (event.target as HTMLTextAreaElement)?.value ?? event.currentTarget?.value ?? "";
    props.onChange?.(val);
  };

  return (
    <div class="flex flex-col gap-1.5">
      <Label for={id}>
        {props.field.label}
        <Show when={props.field.required}>
          <span aria-hidden="true" class="ml-0.5 text-destructive">
            *
          </span>
        </Show>
      </Label>

      <Textarea
        id={id}
        name={props.field.name}
        rows={props.field.rows ?? 3}
        aria-describedby={describedBy()}
        aria-invalid={props.error ? "true" : undefined}
        value={props.value != null ? String(props.value) : ""}
        placeholder={props.field.placeholder}
        disabled={props.field.disabled}
        onInput={handleInput}
      />

      <Show
        when={props.error}
        fallback={
          <Show when={props.field.hint}>
            <p id={hintId} class="text-xs text-muted-foreground">
              {props.field.hint}
            </p>
          </Show>
        }
      >
        {(err) => (
          <p id={errorId} role="alert" class="text-xs font-medium text-destructive">
            {err()}
          </p>
        )}
      </Show>
    </div>
  );
}
