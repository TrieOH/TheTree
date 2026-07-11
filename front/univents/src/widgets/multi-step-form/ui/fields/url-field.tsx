import { useState } from "react";
import type { FieldValues } from "react-hook-form";
import { ExternalLink, Link2 } from "lucide-react";
import type { FieldConfig, FieldFormApi } from "../../model/types";
import { getFieldError } from "../../utils/get-field-error";

export interface UrlFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

function looksLikeUrl(value: unknown): value is string {
  if (typeof value !== "string" || value.trim().length === 0) return false;
  try {
    new URL(value);
    return true;
  } catch {
    return false;
  }
}

function withScheme(value: string): string {
  if (/^[a-z][a-z0-9+.-]*:\/\//i.test(value)) return value;
  return `https://${value}`;
}

export function UrlFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: UrlFieldRendererProps<TFieldValues>) {
  if (field.kind !== "url") return null;

  const error = getFieldError(form.formState.errors, field.name);
  const currentValue = form.watch(field.name);
  const canOpen = looksLikeUrl(currentValue);
  const autoPrependScheme = field.autoPrependScheme ?? true;

  const { onBlur: registerOnBlur, ...registerRest } = form.register(field.name);
  const [isFocused, setIsFocused] = useState(false);

  return (
    <div className="space-y-1.5">
      <label
        htmlFor={field.name}
        className="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
      >
        {field.label}
        {field.optional ? (
          <span className="ml-1 font-normal normal-case text-muted-foreground/70">(opcional)</span>
        ) : null}
      </label>

      <div
        className={
          "flex h-10 items-center gap-2 rounded-md border bg-transparent px-3 text-sm " +
          (error
            ? "border-destructive"
            : isFocused
              ? "border-ring ring-2 ring-ring"
              : "border-input")
        }
      >
        <Link2 className="h-4 w-4 shrink-0 text-muted-foreground" />

        <input
          id={field.name}
          type="text"
          inputMode="url"
          placeholder={field.placeholder ?? "https://..."}
          disabled={field.disabled}
          aria-invalid={Boolean(error)}
          className="h-full w-full bg-transparent outline-none placeholder:text-muted-foreground"
          {...registerRest}
          onFocus={() => setIsFocused(true)}
          onBlur={(event) => {
            setIsFocused(false);
            if (autoPrependScheme && event.target.value.trim().length > 0) {
              const normalized = withScheme(event.target.value.trim());
              if (normalized !== event.target.value) {
                form.setValue(field.name, normalized as never, { shouldValidate: true, shouldDirty: true });
              }
            }
            void registerOnBlur(event);
          }}
        />

        {canOpen ? (
          <a
            href={currentValue}
            target="_blank"
            rel="noreferrer noopener"
            aria-label="Abrir link em nova aba"
            className="shrink-0 text-muted-foreground hover:text-foreground"
          >
            <ExternalLink className="h-4 w-4" />
          </a>
        ) : null}
      </div>

      {field.description ? <p className="text-xs text-muted-foreground">{field.description}</p> : null}
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  );
}
