import { shapeClass, type Shape } from "./lib/shape";
import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal, onSettled } from "solid-js";
import { cn } from "./lib/cn";

export interface InputProps {
  shape?: Shape;
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
  "aria-label"?: string;
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
    props.type === "date" || props.type === "datetime-local";

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
      aria-label={props["aria-label"]}
      aria-describedby={props["aria-describedby"]}
      aria-invalid={ariaInvalid()}
      onInput={props.onInput}
      onChange={props.onChange}
      onBlur={props.onBlur}
      onFocus={props.onFocus}
      onKeyDown={props.onKeyDown}
      class={cn(
        fieldBase,
        isDateOrTime() && "relative pr-9",
        shapeClass(props.shape),
        props.class,
      )}
    />
  );
}

export type TextareaResize = "none" | "vertical" | "horizontal" | "both";

export interface TextareaProps {
  shape?: Shape;
  ref?: HTMLTextAreaElement | ((el: HTMLTextAreaElement) => void);
  id?: string;
  name?: string;
  value?: string;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  readonly?: boolean;
  /** Visible lines, and the minimum height when autosizing. */
  rows?: number;
  /**
   * Hard character limit: the browser stops accepting input at this length.
   * Shows the character counter unless `showCount` is `false`, so the limit is
   * visible before you run into it.
   */
  maxlength?: number;
  /**
   * Live character count under the field. With `maxlength` it reads "12/280"
   * and emphasises itself as the limit gets close; without one it counts
   * characters. On by default when `maxlength` is set, off otherwise.
   */
  showCount?: boolean;
  /** Grows with its content instead of scrolling. */
  autosize?: boolean;
  /** Tallest an autosizing field gets, in lines, before it scrolls. */
  maxRows?: number;
  /**
   * Which way the user can drag to resize. Vertical by default, none while
   * `autosize` is on (dragging would fight the automatic height).
   */
  resize?: TextareaResize;
  class?: string;
  "aria-label"?: string;
  "aria-describedby"?: string;
  "aria-invalid"?: boolean | "true" | "false" | "grammar" | "spelling";
  onInput?: JSX.EventHandler<HTMLTextAreaElement, InputEvent>;
  onChange?: JSX.EventHandler<HTMLTextAreaElement, Event>;
  onBlur?: JSX.EventHandler<HTMLTextAreaElement, FocusEvent>;
  onFocus?: JSX.EventHandler<HTMLTextAreaElement, FocusEvent>;
  onKeyDown?: JSX.EventHandler<HTMLTextAreaElement, KeyboardEvent>;
}

const resizeClass: Record<TextareaResize, string> = {
  none: "resize-none",
  vertical: "resize-y",
  horizontal: "resize-x",
  both: "resize",
};

export function Textarea(props: TextareaProps) {
  const ariaInvalid = () => {
    const val = props["aria-invalid"];
    if (val === true || val === "true") return "true";
    if (val === "grammar" || val === "spelling") return val;
    return undefined;
  };

  const counterId = `textarea-count-${Math.random().toString(36).slice(2, 9)}`;
  let element: HTMLTextAreaElement | undefined;
  // Length of what the user typed, for uncontrolled use; a `value` prop wins.
  const [typed, setTyped] = createSignal(0);
  const count = () => (props.value !== undefined ? props.value.length : typed());
  const resize = () => props.resize ?? (props.autosize ? "none" : "vertical");
  const showCount = () => props.showCount ?? props.maxlength !== undefined;

  const describedBy = () =>
    [props["aria-describedby"], showCount() ? counterId : undefined]
      .filter(Boolean)
      .join(" ") || undefined;

  const fit = () => {
    const el = element;
    if (!el || !props.autosize) return;

    el.style.height = "auto";
    const style = getComputedStyle(el);
    const borders =
      (parseFloat(style.borderTopWidth) || 0) +
      (parseFloat(style.borderBottomWidth) || 0);
    const padding =
      (parseFloat(style.paddingTop) || 0) + (parseFloat(style.paddingBottom) || 0);
    const line =
      parseFloat(style.lineHeight) || (parseFloat(style.fontSize) || 14) * 1.5;
    const content = el.scrollHeight + borders;
    const limit =
      props.maxRows !== undefined
        ? line * props.maxRows + padding + borders
        : Infinity;

    el.style.height = `${Math.min(content, limit)}px`;
    el.style.overflowY = content > limit ? "auto" : "hidden";
  };

  // Re-fit when the text, the mode or the cap changes from outside.
  createEffect(
    () => [props.value, props.autosize, props.maxRows] as const,
    () => fit(),
  );

  // Wrapping depends on width, so a width change needs a re-fit too.
  onSettled(() => {
    if (!element || typeof ResizeObserver === "undefined") return;
    let width = element.offsetWidth;
    const observer = new ResizeObserver(() => {
      if (!element || element.offsetWidth === width) return;
      width = element.offsetWidth;
      fit();
    });
    observer.observe(element);
    return () => observer.disconnect();
  });

  const control = (
    <textarea
      ref={(el: HTMLTextAreaElement) => {
        element = el;
        const forward = props.ref;
        if (typeof forward === "function") forward(el);
      }}
      id={props.id}
      name={props.name}
      value={props.value ?? ""}
      placeholder={props.placeholder}
      disabled={props.disabled}
      required={props.required}
      readonly={props.readonly}
      rows={props.rows ?? 3}
      maxlength={props.maxlength}
      aria-label={props["aria-label"]}
      aria-describedby={describedBy()}
      aria-invalid={ariaInvalid()}
      onInput={(event) => {
        setTyped(event.currentTarget.value.length);
        fit();
        props.onInput?.(event);
      }}
      onChange={props.onChange}
      onBlur={props.onBlur}
      onFocus={props.onFocus}
      onKeyDown={props.onKeyDown}
      class={cn(
        fieldBase,
        "min-h-20 custom-scrollbar",
        resizeClass[resize()],
        shapeClass(props.shape, "surface"),
        props.class,
      )}
    />
  );

  return (
    <Show when={showCount()} fallback={control}>
      <div class="flex w-full flex-col gap-1">
        {control}
        <p
          id={counterId}
          class={cn(
            "text-right text-xs tabular-nums text-muted-foreground",
            props.maxlength !== undefined &&
              count() >= props.maxlength * 0.9 &&
              "text-foreground",
            props.maxlength !== undefined &&
              count() >= props.maxlength &&
              "font-medium text-destructive dark:text-destructive-foreground",
          )}
        >
          {props.maxlength !== undefined
            ? `${count()}/${props.maxlength}`
            : `${count()} caracteres`}
        </p>
      </div>
    </Show>
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
      <p
        id={props.id}
        role="alert"
        class="border-l-2 border-destructive dark:border-destructive-foreground pl-2 text-xs text-foreground"
      >
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
          <span
            aria-hidden="true"
            class="ml-0.5 text-destructive dark:text-destructive-foreground"
          >
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
        <p
          id={errorId()}
          role="alert"
          class="border-l-2 border-destructive dark:border-destructive-foreground pl-2 text-xs font-medium text-foreground"
        >
          {props.error}
        </p>
      </Show>
    </div>
  );
}
