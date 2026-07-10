import type { FieldValues } from "react-hook-form";
import type { FieldConfig, FieldFormApi } from "../../model/types";
import { getFieldError } from "../../utils/get-field-error";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";

export interface TextFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

export function TextFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: TextFieldRendererProps<TFieldValues>) {
  if (field.kind !== "text") return null;

  const error = getFieldError(form.formState.errors, field.name);

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
      <Input
        id={field.name}
        type={field.inputType ?? "text"}
        placeholder={field.placeholder}
        disabled={field.disabled}
        aria-invalid={Boolean(error)}
        className={error ? "border-destructive focus-visible:ring-destructive" : undefined}
        {...form.register(field.name)}
      />
      {field.description ? <p className="text-xs text-muted-foreground">{field.description}</p> : null}
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  );
}
