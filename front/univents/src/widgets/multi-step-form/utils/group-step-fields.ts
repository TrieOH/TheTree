import type { FieldValues } from "react-hook-form";
import type { FieldConfig } from "../model/types";

export function groupStepFields<TFieldValues extends FieldValues>(
  fields: FieldConfig<TFieldValues>[],
): FieldConfig<TFieldValues>[][] {
  const rows: FieldConfig<TFieldValues>[][] = [];
  let pendingHalfField: FieldConfig<TFieldValues> | null = null;

  for (const field of fields) {
    if (field.layout === "half") {
      if (pendingHalfField) {
        rows.push([pendingHalfField, field]);
        pendingHalfField = null;
      } else pendingHalfField = field;
      continue;
    }

    if (pendingHalfField) {
      rows.push([pendingHalfField]);
      pendingHalfField = null;
    }

    rows.push([field]);
  }

  if (pendingHalfField) rows.push([pendingHalfField]);

  return rows;
}
