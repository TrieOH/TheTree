import type { FieldI } from "#/features/fields/model";
import { cn } from "#/shared/lib/utils";
import type { FieldAnswerable, StepAnswerable } from "@trieoh/informd-models";

interface ReviewStepProps {
  steps: StepAnswerable[];
  fields: Record<string, FieldAnswerable[]>;
  formData: Record<string, unknown>;
}

function formatValue(field: FieldI, value: unknown): string {
  if (value === undefined || value === "" || value === null) return "-";

  if (field.type === "bool") {
    return value === true || value === "true" ? "Yes" : "No";
  }


  if (field.type === "select") {
    const options = field.config?.options ?? [];
    if (Array.isArray(value)) {
      return value
        .map((v) => options.find((o: any) => o.value === v)?.label ?? v)
        .join(", ");
    }
    return options.find((o: any) => o.value === value)?.label ?? String(value);
  }

  if (field.type === "file") {
    return (value as File).name;
  }

  return String(value);
}

export function ReviewStep({ steps, fields, formData }: ReviewStepProps) {
  // Only review steps that have fields
  const reviewSteps = steps.filter((s) => fields[s.step.id].length > 0);

  return (
    <div className="space-y-4">
      <div className="rounded-lg bg-muted/30 p-4">
        {reviewSteps.map((step, stepIdx) => {
          const stepFields = fields[step.step.id].sort(
            (a, b) => a.field.position_hint - b.field.position_hint
          );

          return (
            <div
              key={step.step.id}
              className={cn(
                stepIdx < reviewSteps.length - 1 && "mb-5 border-b border-border/50 pb-5"
              )}
            >
              <h3 className="mb-2 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                {step.step.title}
              </h3>
              <div className="space-y-2">
                {stepFields.map((field) => (
                  <div
                    key={field.field.id}
                    className="flex justify-between gap-4 border-b border-border/30 py-1.5 last:border-0"
                  >
                    <span className="text-xs text-muted-foreground">{field.field.title}</span>
                    <span className="max-w-[60%] text-right text-xs font-medium text-foreground">
                      {formatValue(field.field, formData[field.field.id])}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}