import type { FieldErrors, FieldValues } from "react-hook-form";

/**
 * RHF's FieldErrors type is deeply nested and keyed by the shape of
 * TFieldValues, which makes indexing it with a dynamic dot-path string
 * awkward to do without `any`. This walks the path manually and narrows
 * safely at each step.
 */
export function getFieldError<TFieldValues extends FieldValues>(
  errors: FieldErrors<TFieldValues>,
  path: string,
): string | undefined {
  let current: unknown = errors;

  for (const segment of path.split(".")) {
    if (current && typeof current === "object" && segment in current)
      current = (current as Record<string, unknown>)[segment];
    else return undefined;
  }

  if (current && typeof current === "object" && "message" in current) {
    const message = (current as { message?: unknown }).message;
    return typeof message === "string" ? message : undefined;
  }

  return undefined;
}
