import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

const fieldBase =
  "w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive";

export interface InputProps {
  id?: string;
  name?: string;
  type?: "text" | "email" | "password" | "search" | "url" | "tel" | "number" | "date";
  value?: string | number;
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  readonly?: boolean;
  autocomplete?: string;
  inputmode?: "none" | "text" | "decimal" | "numeric" | "tel" | "email" | "url" | "search";
  min?: string | number;
  max?: string | number;
  step?: string | number;
  minlength?: number;
  maxlength?: number;
  class?: string;
  /** Solid 2 types ARIA booleans as the string form. */
  "aria-invalid"?: "true" | "false";
  "aria-describedby"?: string;
  onInput?: (event: InputEvent & { currentTarget: HTMLInputElement }) => void;
  onChange?: (event: Event & { currentTarget: HTMLInputElement }) => void;
  onBlur?: (event: FocusEvent & { currentTarget: HTMLInputElement }) => void;
}

export function Input(props: InputProps) {
  return (
    <input
      id={props.id}
      name={props.name}
      type={props.type ?? "text"}
      value={props.value}
      placeholder={props.placeholder}
      required={props.required}
      disabled={props.disabled}
      readonly={props.readonly}
      autocomplete={props.autocomplete}
      inputmode={props.inputmode}
      min={props.min}
      max={props.max}
      step={props.step}
      minlength={props.minlength}
      maxlength={props.maxlength}
      aria-invalid={props["aria-invalid"]}
      aria-describedby={props["aria-describedby"]}
      onInput={props.onInput}
      onChange={props.onChange}
      onBlur={props.onBlur}
      class={cn(fieldBase, "h-9", props.class)}
    />
  );
}

export interface TextareaProps {
  id?: string;
  name?: string;
  value?: string;
  placeholder?: string;
  rows?: number;
  required?: boolean;
  disabled?: boolean;
  readonly?: boolean;
  class?: string;
  "aria-invalid"?: "true" | "false";
  "aria-describedby"?: string;
  onInput?: (event: InputEvent & { currentTarget: HTMLTextAreaElement }) => void;
  onBlur?: (event: FocusEvent & { currentTarget: HTMLTextAreaElement }) => void;
}

export function Textarea(props: TextareaProps) {
  return (
    <textarea
      id={props.id}
      name={props.name}
      value={props.value}
      placeholder={props.placeholder}
      rows={props.rows}
      required={props.required}
      disabled={props.disabled}
      readonly={props.readonly}
      aria-invalid={props["aria-invalid"]}
      aria-describedby={props["aria-describedby"]}
      onInput={props.onInput}
      onBlur={props.onBlur}
      class={cn(fieldBase, "min-h-20", props.class)}
    />
  );
}

export interface LabelProps {
  for?: string;
  class?: string;
  children?: JSX.Element;
}

export function Label(props: LabelProps) {
  return (
    <label
      for={props.for}
      class={cn("text-sm font-medium text-foreground", props.class)}
    >
      {props.children}
    </label>
  );
}

export interface FieldProps {
  label: JSX.Element;
  /** Rendered below the control, and wired to `aria-describedby`. */
  hint?: JSX.Element;
  error?: string;
  required?: boolean;
  /** Receives the generated id so the caller can bind control and label. */
  children: (ids: { id: string; describedBy?: string }) => JSX.Element;
  class?: string;
}

/**
 * Label + control + message, with the `aria` wiring in one place: hand-rolled
 * so a form cannot ship a label that is not associated with its input, which is
 * the part everyone forgets when wiring fields by hand.
 */
export function Field(props: FieldProps) {
  const id = `field-${Math.random().toString(36).slice(2, 9)}`;
  const errorId = `${id}-error`;
  const hintId = `${id}-hint`;
  const describedBy = () => {
    if (props.error) return errorId;
    if (props.hint) return hintId;
    return undefined;
  };

  return (
    <div class={cn("flex flex-col gap-1.5", props.class)}>
      <Label for={id}>
        {props.label}
        <Show when={props.required}>
          <span aria-hidden="true" class="ml-0.5 text-destructive">
            *
          </span>
        </Show>
      </Label>

      {props.children({ id, describedBy: describedBy() })}

      <Show when={props.error} fallback={
        <Show when={props.hint}>
          <p id={hintId} class="text-xs text-muted-foreground">
            {props.hint}
          </p>
        </Show>
      }>
        <p id={errorId} role="alert" class="text-xs text-destructive">
          {props.error}
        </p>
      </Show>
    </div>
  );
}
