import type { JSX } from "@solidjs/web";
import { createMemo, Show } from "solid-js";
import { Input, Label } from "../../input";
import type { MoneyMultiStepField } from "../types";

export interface MultiStepMoneyFieldProps<T = Record<string, unknown>> {
  field: MoneyMultiStepField<T>;
  value?: unknown;
  error?: string;
  onChange?: (value: number) => void;
}

function toCents(val: unknown): number {
  if (typeof val === "number" && Number.isFinite(val)) return Math.round(val);
  if (typeof val === "string") {
    const digits = val.replace(/\D/g, "");
    return digits ? parseInt(digits, 10) : 0;
  }
  return 0;
}

export function MultiStepMoneyField<T = Record<string, unknown>>(
  props: MultiStepMoneyFieldProps<T>,
): JSX.Element {
  let inputRef: HTMLInputElement | undefined;

  const id = `field-${props.field.name}`;
  const errorId = `${id}-error`;
  const hintId = `${id}-hint`;

  const describedBy = () => {
    if (props.error) return errorId;
    if (props.field.hint) return hintId;
    return undefined;
  };

  const currency = () => props.field.currency ?? "BRL";
  const locale = () => props.field.locale ?? "pt-BR";

  const cents = () => toCents(props.value);

  const formatter = createMemo(() => {
    return new Intl.NumberFormat(locale(), {
      style: "currency",
      currency: currency(),
    });
  });

  const display = () => formatter().format(cents() / 100);

  const commit = (nextCents: number) => {
    let result = nextCents;
    if (typeof props.field.minCents === "number") {
      result = Math.max(result, props.field.minCents);
    }
    if (typeof props.field.maxCents === "number") {
      result = Math.min(result, props.field.maxCents);
    }
    props.onChange?.(result);
  };

  const handleInput = (
    event: InputEvent & { currentTarget: HTMLInputElement },
  ) => {
    const target =
      (event.target as HTMLInputElement) ??
      (event.currentTarget as HTMLInputElement);
    const raw = (target?.value ?? "").replace(/\D/g, "");
    const parsed = raw ? parseInt(raw, 10) : 0;
    commit(parsed);
  };

  const handleKeyDown = (
    event: KeyboardEvent & { currentTarget: HTMLInputElement },
  ) => {
    if (event.key === "Backspace") {
      event.preventDefault();
      const current = cents();
      const nextCents = Math.trunc(current / 10);
      commit(nextCents);
    }
  };

  const handleFocus = () => {
    requestAnimationFrame(() => {
      const len = inputRef?.value.length ?? 0;
      inputRef?.setSelectionRange(len, len);
    });
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
        ref={(el) => (inputRef = el)}
        id={id}
        name={props.field.name}
        type="text"
        inputmode="numeric"
        aria-describedby={describedBy()}
        aria-invalid={props.error ? "true" : undefined}
        value={display()}
        disabled={props.field.disabled}
        autocomplete="off"
        onInput={handleInput}
        onKeyDown={handleKeyDown}
        onFocus={handleFocus}
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
