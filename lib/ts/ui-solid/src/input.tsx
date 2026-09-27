import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

export interface InputProps {
  ref?: HTMLInputElement | ((el: HTMLInputElement) => void);
  id?: string;
  name?: string;
  type?:
    | "text"
    | "email"
    | "password"
    | "search"
    | "url"
    | "tel"
    | "number"
    | "date"
    | "datetime-local";
  value?: string | number;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  readonly?: boolean;
  autocomplete?: string;
  inputmode?:
    | "none"
    | "text"
    | "decimal"
    | "numeric"
    | "tel"
    | "email"
    | "url"
    | "search";
  min?: string | number;
  max?: string | number;
  step?: string | number;
  minlength?: number;
  maxlength?: number;
  class?: string;
  "aria-describedby"?: string;
  "aria-invalid"?: boolean | "true" | "false" | "grammar" | "spelling";
  onInput?: JSX.EventHandler<HTMLInputElement, InputEvent>;
  onChange?: JSX.EventHandler<HTMLInputElement, Event>;
  onBlur?: JSX.EventHandler<HTMLInputElement, FocusEvent>;
  onFocus?: JSX.EventHandler<HTMLInputElement, FocusEvent>;
  onKeyDown?: JSX.EventHandler<HTMLInputElement, KeyboardEvent>;
}

const fieldBase =
  "flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50";

export function Input(props: InputProps) {
  const ariaInvalid = () => {
    const val = props["aria-invalid"];
    if (val === true || val === "true") return "true";
    if (val === "grammar" || val === "spelling") return val;
    return undefined;
  };

  const isDateOrTime = () =>
    props.type === "date" ||
    props.type === "datetime-local";

  return (
    <input
      ref={props.ref}
      id={props.id}
      name={props.name}
      type={props.type ?? "text"}
      value={props.value ?? ""}
      placeholder={props.placeholder}
      disabled={props.disabled}
      required={props.required}
      readonly={props.readonly}
      autocomplete={props.autocomplete}
      inputmode={props.inputmode}
      min={props.min}
      max={props.max}
      step={props.step}
      minlength={props.minlength}
      maxlength={props.maxlength}
      aria-describedby={props["aria-describedby"]}
      aria-invalid={ariaInvalid()}
      onInput={props.onInput}
      onChange={props.onChange}
      onBlur={props.onBlur}
      onFocus={props.onFocus}
      onKeyDown={props.onKeyDown}
      class={cn(fieldBase, isDateOrTime() && "relative pr-9", props.class)}
    />
  );
}

export interface TextareaProps {
  ref?: HTMLTextAreaElement | ((el: HTMLTextAreaElement) => void);
  id?: string;
  name?: string;
  value?: string;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  readonly?: boolean;
  rows?: number;
  class?: string;
  "aria-describedby"?: string;
  "aria-invalid"?: boolean | "true" | "false" | "grammar" | "spelling";
  onInput?: JSX.EventHandler<HTMLTextAreaElement, InputEvent>;
  onChange?: JSX.EventHandler<HTMLTextAreaElement, Event>;
  onBlur?: JSX.EventHandler<HTMLTextAreaElement, FocusEvent>;
  onFocus?: JSX.EventHandler<HTMLTextAreaElement, FocusEvent>;
  onKeyDown?: JSX.EventHandler<HTMLTextAreaElement, KeyboardEvent>;
}

export function Textarea(props: TextareaProps) {
  const ariaInvalid = () => {
    const val = props["aria-invalid"];
    if (val === true || val === "true") return "true";
    if (val === "grammar" || val === "spelling") return val;
    return undefined;
  };

  return (
    <textarea
      ref={props.ref}
      id={props.id}
      name={props.name}
      value={props.value ?? ""}
      placeholder={props.placeholder}
      disabled={props.disabled}
      required={props.required}
      readonly={props.readonly}
      rows={props.rows ?? 3}
      aria-describedby={props["aria-describedby"]}
      aria-invalid={ariaInvalid()}
      onInput={props.onInput}
      onChange={props.onChange}
      onBlur={props.onBlur}
      onFocus={props.onFocus}
      onKeyDown={props.onKeyDown}
      class={cn(fieldBase, "min-h-20 custom-scrollbar", props.class)}
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

export interface FormErrorProps {
  id?: string;
  children?: JSX.Element;
}

export function FormError(props: FormErrorProps) {
  return (
    <Show when={props.children}>
      <p id={props.id} role="alert" class="text-xs text-destructive">
        {props.children}
      </p>
    </Show>
  );
}

export interface FieldProps {
  id?: string;
  name?: string;
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
  const autoId = `field-${Math.random().toString(36).slice(2, 9)}`;
  const id = () => props.id ?? (props.name ? `field-${props.name}` : autoId);
  const errorId = () => `${id()}-error`;
  const hintId = () => `${id()}-hint`;

  const ids = {
    get id() {
      return id();
    },
    get describedBy() {
      if (props.error) return errorId();
      if (props.hint) return hintId();
      return undefined;
    },
  };

  return (
    <div class={cn("flex flex-col gap-1.5", props.class)}>
      <Label for={id()}>
        {props.label}
        <Show when={props.required}>
          <span aria-hidden="true" class="ml-0.5 text-destructive">
            *
          </span>
        </Show>
      </Label>

      {props.children(ids)}

      <Show
        when={props.error}
        fallback={
          <Show when={props.hint}>
            <p id={hintId()} class="text-xs text-muted-foreground">
              {props.hint}
            </p>
          </Show>
        }
      >
        <p id={errorId()} role="alert" class="text-xs font-medium text-destructive">
          {props.error}
        </p>
      </Show>
    </div>
  );
}
