import { useMemo, useRef, useState } from "react";
import type { FieldValues } from "react-hook-form";
import type { FieldConfig, FieldFormApi } from "../../model/types";
import { getFieldError } from "../../utils/get-field-error";

export interface MoneyFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

function toCentsNumber(value: unknown): number {
  if (typeof value === "bigint") return Number(value);
  if (typeof value === "number" && Number.isFinite(value)) return value;
  return 0;
}

function clamp(cents: number, min?: number, max?: number): number {
  let result = cents;
  if (typeof min === "number") result = Math.max(result, min);
  if (typeof max === "number") result = Math.min(result, max);
  return result;
}

/** Keeps only the digits typed so far — this is what makes new digits
 * "shift in from the right" instead of being parsed positionally. */
function digitsToCents(raw: string): number {
  const digitsOnly = raw.replace(/\D/g, "");
  if (digitsOnly.length === 0) return 0;
  return parseInt(digitsOnly, 10);
}

export function MoneyFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: MoneyFieldRendererProps<TFieldValues>) {
  if (field.kind !== "money") return null;

  const inputRef = useRef<HTMLInputElement>(null);
  const error = getFieldError(form.formState.errors, field.name);
  const currency = field.currency ?? "BRL";
  const locale = field.locale ?? "pt-BR";
  const valueType = field.valueType ?? "number";

  const [cents, setCents] = useState<number>(() =>
    toCentsNumber(form.getValues(field.name)),
  );

  const formatter = useMemo(
    () => new Intl.NumberFormat(locale, { style: "currency", currency }),
    [locale, currency],
  );

  const display = formatter.format(cents / 100);

  const commit = (nextCents: number) => {
    const clamped = clamp(nextCents, field.minCents, field.maxCents);
    setCents(clamped);
    const nextValue = valueType === "bigint" ? BigInt(clamped) : clamped;
    form.setValue(field.name, nextValue as never, {
      shouldDirty: true,
      shouldValidate: true,
    });
  };

  return (
    <div className="space-y-1.5">
      <label
        htmlFor={field.name}
        className="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
      >
        {field.label}
        {field.optional ? (
          <span className="ml-1 font-normal normal-case text-muted-foreground/70">
            (opcional)
          </span>
        ) : null}
      </label>

      <input
        ref={inputRef}
        id={field.name}
        type="text"
        inputMode="numeric"
        disabled={field.disabled}
        aria-invalid={Boolean(error)}
        value={display}
        onChange={(event) => commit(digitsToCents(event.target.value))}
        onKeyDown={(event) => {
          // Backspace should drop the last digit even though the field
          // is fully controlled/re-formatted on every keystroke.
          if (event.key === "Backspace") {
            event.preventDefault();
            commit(Math.trunc(cents / 10));
          }
        }}
        onFocus={() => {
          requestAnimationFrame(() => {
            const length = inputRef.current?.value.length ?? 0;
            inputRef.current?.setSelectionRange(length, length);
          });
        }}
        className={
          "flex h-10 w-full rounded-md border bg-transparent px-3 py-2 text-left text-sm outline-none " +
          "tabular-nums focus-visible:ring-2 focus-visible:ring-ring " +
          (error
            ? "border-destructive focus-visible:ring-destructive"
            : "border-input")
        }
      />

      {field.description ? (
        <p className="text-xs text-muted-foreground">{field.description}</p>
      ) : null}
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  );
}
