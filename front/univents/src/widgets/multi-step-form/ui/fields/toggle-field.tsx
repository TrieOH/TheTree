import type { FieldValues } from "react-hook-form";
import type { FieldConfig, FieldFormApi } from "../../model/types";

export interface ToggleFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

export function ToggleFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: ToggleFieldRendererProps<TFieldValues>) {
  if (field.kind !== "toggle") return null;

  const checked = Boolean(form.watch(field.name));

  return (
    <label
      htmlFor={field.name}
      className={
        "flex items-center justify-between gap-4 rounded-md border px-3 py-2.5 " +
        (field.disabled ? "opacity-60" : "cursor-pointer")
      }
    >
      <span className="space-y-0.5">
        <span className="block text-sm font-medium">{field.label}</span>
        {field.description ? (
          <span className="block text-xs text-muted-foreground">
            {field.description}
          </span>
        ) : null}
      </span>

      <span className="relative inline-flex h-6 w-11 shrink-0 items-center">
        <input
          id={field.name}
          type="checkbox"
          disabled={field.disabled}
          className="peer sr-only"
          {...form.register(field.name)}
        />
        <span
          aria-hidden="true"
          className={
            "absolute inset-0 rounded-full transition-colors " +
            (checked ? "bg-foreground" : "bg-input") +
            " peer-focus-visible:ring-2 peer-focus-visible:ring-ring peer-focus-visible:ring-offset-2"
          }
        />
        <span
          aria-hidden="true"
          className={
            "relative h-5 w-5 rounded-full bg-background shadow transition-transform " +
            (checked ? "translate-x-5" : "translate-x-0.5")
          }
        />
      </span>
    </label>
  );
}
