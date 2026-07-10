import type { FieldValues } from "react-hook-form";
import { ArrowLeft, ArrowRight } from "lucide-react";
import { createFieldRegistry, renderField } from "./field-registry";
import type { MultiStepFormController } from "../hooks/use-multi-step-form";
import { Button } from "@/shared/ui/shadcn/button";

export interface MultiStepFormProps<TInput extends FieldValues, TOutput = TInput> {
  controller: MultiStepFormController<TInput, TOutput>;
  submitLabel?: string;
  onCancel?: () => void;
}

export function MultiStepForm<TInput extends FieldValues, TOutput = TInput>({
  controller,
  submitLabel = "Concluir",
  onCancel,
}: MultiStepFormProps<TInput, TOutput>) {
  const {
    form,
    steps,
    stepIndex,
    isFirstStep,
    isLastStep,
    visibleFields,
    goNext,
    goBack,
    isSubmitting,
    canSubmit,
  } = controller;

  const registry = createFieldRegistry<TInput>();

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <ol className="mb-6 flex flex-wrap items-center gap-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
        {steps.map((step, index) => (
          <li key={step.id} className="flex items-center gap-2">
            <span className="flex items-center gap-1.5">
              <span
                className={
                  "flex h-5 w-6 items-center justify-center rounded border text-[10px] " +
                  (index === stepIndex
                    ? "border-foreground bg-foreground text-background"
                    : "border-border text-muted-foreground")
                }
              >
                {String(index + 1).padStart(2, "0")}
              </span>
              <span className={index === stepIndex ? "text-foreground" : undefined}>{step.label}</span>
            </span>
            {index < steps.length - 1 ? <span className="text-border">{"\u203A"}</span> : null}
          </li>
        ))}
      </ol>

      <form
        onSubmit={(event) => {
          event.preventDefault();
          void goNext();
        }}
        className="flex min-h-0 flex-1 flex-col"
      >
        <div className="min-h-0 flex-1 space-y-5 overflow-y-auto px-1">
          {visibleFields.map((field) => renderField(field, form, registry))}
        </div>

        <div className="mt-6 flex shrink-0 items-center justify-between border-t pt-4">
          {!isFirstStep ? (
            <Button
              type="button"
              onClick={goBack}
              disabled={isSubmitting}
              variant="ghost"
              className="text-sm font-medium text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="size-4" />
              Voltar
            </Button>
          ) : onCancel ? (
            <Button
              type="button"
              onClick={onCancel}
              disabled={isSubmitting}
              variant="ghost"
              className="text-sm font-medium text-muted-foreground hover:text-foreground"
            >
              Cancelar
            </Button>
          ) : (
            <span />
          )}

          <Button
            type="submit"
            disabled={isSubmitting || (isLastStep && !canSubmit)}
            title={isLastStep && !canSubmit ? "Nenhuma alteração para salvar" : undefined}
            className="p-4 text-sm font-semibold"
          >
            {isLastStep ? submitLabel : "Avançar"}
            {!isLastStep ? <ArrowRight className="size-4" /> : null}
          </Button>
        </div>
      </form>
    </div>
  );
}
