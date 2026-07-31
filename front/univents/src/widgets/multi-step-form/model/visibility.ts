import type { FieldValues } from "react-hook-form";
import type { VisibilityRule } from "./types";

export function isFieldVisible<TFieldValues extends FieldValues>(
  rule: VisibilityRule<TFieldValues> | undefined,
  values: TFieldValues,
): boolean {
  if (!rule) return true;

  switch (rule.type) {
    case "equals":
      return getPathValue(values, rule.field) === rule.value;
    case "notEquals":
      return getPathValue(values, rule.field) !== rule.value;
    case "hasValue":
      return hasMeaningfulValue(getPathValue(values, rule.field));
    case "predicate":
      return rule.predicate(values);
  }
}

function hasMeaningfulValue(value: unknown): boolean {
  if (value === null || value === undefined) return false;
  if (typeof value === "string") return value.trim().length > 0;
  if (Array.isArray(value)) return value.length > 0;
  return true;
}

function getPathValue(source: unknown, path: string): unknown {
  return path.split(".").reduce<unknown>((acc, key) => {
    if (
      acc &&
      typeof acc === "object" &&
      key in (acc as Record<string, unknown>)
    )
      return (acc as Record<string, unknown>)[key];
    return undefined;
  }, source);
}
