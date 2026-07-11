import type { FieldValues } from "react-hook-form";
import type { FieldConfig, FieldFormApi } from "../../model/types";

export interface CustomFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

export function CustomFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: CustomFieldRendererProps<TFieldValues>) {
  if (field.kind !== "custom") return null;
  return <>{field.render({ form })}</>;
}