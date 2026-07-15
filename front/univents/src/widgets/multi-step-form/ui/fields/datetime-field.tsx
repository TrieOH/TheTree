import type { FieldValues } from "react-hook-form";
import type { FieldConfig, FieldFormApi } from "../../model/types";
import { getFieldError } from "../../utils/get-field-error";
import { DateTimePicker } from "@/widgets/multi-step-form/ui/fields/helper/date-time-picker";
import { Label } from "@/shared/ui/shadcn/label";

export interface DateTimeFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

export function DateTimeFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: DateTimeFieldRendererProps<TFieldValues>) {
  if (field.kind !== "datetime") return null;

  const error = getFieldError(form.formState.errors, field.name);
  const currentValue = form.watch(field.name);

  return (
    <div className="space-y-1.5">
      <Label
        htmlFor={field.name}
        className="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
      >
        {field.label}
        {field.optional ? (
          <span className="ml-1 font-normal normal-case text-muted-foreground/70">(opcional)</span>
        ) : null}
      </Label>

      <DateTimePicker
        id={field.name}
        value={typeof currentValue === "string" ? currentValue : ""}
        onChange={(value) => {
          form.setValue(field.name, value as never, {
            shouldDirty: true,
            shouldValidate: true,
          });
        }}
        disabled={field.disabled}
        placeholder={field.placeholder}
        min={field.min}
        max={field.max}
        error={Boolean(error)}
      />

      {field.description ? <p className="text-xs text-muted-foreground">{field.description}</p> : null}
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  );
}
