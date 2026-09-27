import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { Input, Label } from "../../input";
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
  const id = `field-${props.field.name}`;
  const errorId = `${id}-error`;
  const hintId = `${id}-hint`;

  const describedBy = () => {
    if (props.error) return errorId;
    if (props.field.hint) return hintId;
    return undefined;
  };

  const handleInput = (event: InputEvent & { currentTarget: HTMLInputElement }) => {
    const target = (event.target as HTMLInputElement) ?? (event.currentTarget as HTMLInputElement);
    const val = target?.value ?? "";
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

      <Input
        id={id}
        name={props.field.name}
        type={inputType()}
        aria-describedby={describedBy()}
        aria-invalid={props.error ? "true" : undefined}
        value={props.value != null ? String(props.value) : ""}
        placeholder={props.field.placeholder}
        disabled={props.field.disabled}
        autocomplete={props.field.autocomplete ?? "off"}
        min={props.field.min}
        max={props.field.max}
        step={props.field.step}
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
        <p id={errorId} role="alert" class="text-xs font-medium text-destructive">
          {props.error}
        </p>
      </Show>
    </div>
  );
}
