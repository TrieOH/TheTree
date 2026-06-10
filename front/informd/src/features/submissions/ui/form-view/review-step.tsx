import type { FieldI } from "#/features/fields/model";
import type { StepI } from "#/features/steps/model";
import { cn } from "#/shared/lib/utils";

interface ReviewStepProps {
  steps: StepI[];
  fields: Record<string, FieldI[]>;
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
  // Exclui o último step (review) da listagem
  const reviewSteps = steps.slice(0, -1);

  return (
    <div className="space-y-4">
      <div className="rounded-lg bg-muted/30 p-4">
        {reviewSteps.map((step, stepIdx) => {
          const stepFields = fields[step.id].sort(
            (a, b) => a.position_hint - b.position_hint
          );

          return (
            <div
              key={step.id}
              className={cn(
                stepIdx < reviewSteps.length - 1 && "mb-5 border-b border-border/50 pb-5"
              )}
            >
              <h3 className="mb-2 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                {step.title}
              </h3>
              <div className="space-y-2">
                {stepFields.map((field) => (
                  <div
                    key={field.id}
                    className="flex justify-between gap-4 border-b border-border/30 py-1.5 last:border-0"
                  >
                    <span className="text-xs text-muted-foreground">{field.title}</span>
                    <span className="max-w-[60%] text-right text-xs font-medium text-foreground">
                      {formatValue(field, formData[field.id])}
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