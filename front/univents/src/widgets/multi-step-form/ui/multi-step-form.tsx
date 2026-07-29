import { ArrowLeft, ArrowRight, ChevronRight, Loader2 } from "lucide-react";
import { useEffect } from "react";
import type { FieldValues } from "react-hook-form";
import { Button } from "@/shared/ui/shadcn/button";
import { ImageUploadStateProvider } from "../contexts/image-upload-state-context";
import { clearImageUploadTasks } from "../hooks/use-image-upload-queue";
import type { MultiStepFormController } from "../hooks/use-multi-step-form";
import { groupStepFields } from "../utils/group-step-fields";
import { createFieldRegistry, renderField } from "./field-registry";

export interface MultiStepFormProps<
  TInput extends FieldValues,
  TOutput = TInput,
> {
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
    isProcessingUploads,
    canSubmit,
  } = controller;
  const isBusy = isSubmitting || isProcessingUploads;

  const registry = createFieldRegistry<TInput>();
  const fieldRows = groupStepFields(visibleFields);

  useEffect(() => {
    return () => {
      clearImageUploadTasks();
    };
  }, []);

  return (
    <ImageUploadStateProvider>
      <div className="flex h-full min-h-0 flex-1 flex-col">
        <ol className="mb-6 flex flex-wrap items-center gap-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
          {steps.map((step, index) => (
            <li key={step.id} className="flex items-center gap-2">
              <span className="flex items-center gap-1.5">
                <span
                  className={
                    "flex size-6 aspect-square items-center justify-center rounded-md border text-[10px] leading-none " +
                    (index === stepIndex
                      ? "border-foreground bg-foreground text-background"
                      : "border-border text-muted-foreground")
                  }
                >
                  {String(index + 1).padStart(2, "0")}
                </span>
                <span
                  className={
                    index === stepIndex
                      ? "hidden text-foreground sm:inline"
                      : "hidden sm:inline"
                  }
                >
                  {step.label}
                </span>
              </span>
              {index < steps.length - 1 ? (
                <ChevronRight className="size-3 text-border" />
              ) : null}
            </li>
          ))}
        </ol>

        <form
          onSubmit={(event) => {
            event.preventDefault();
            void goNext();
          }}
          aria-busy={isBusy}
          className="flex min-h-0 flex-1 flex-col"
        >
          <div className="relative min-h-0 flex-1 overflow-hidden">
            <div
              className={
                "h-full min-h-0 space-y-4 overflow-y-auto px-1 pr-2 " +
                (isProcessingUploads
                  ? "pointer-events-none select-none opacity-60"
                  : "")
              }
            >
              {fieldRows.map((row) => (
                <div
                  key={row.map((field) => field.name).join("-")}
                  className={
                    row.length > 1
                      ? "grid grid-cols-1 gap-4 md:grid-cols-2"
                      : undefined
                  }
                >
                  {row.map((field) => (
                    <div key={field.name}>
                      {renderField(field, form, registry)}
                    </div>
                  ))}
                </div>
              ))}
            </div>

            {isProcessingUploads ? (
              <div
                className="pointer-events-none absolute inset-0 z-10 rounded-xl bg-background/35 backdrop-blur-[1px]"
                aria-hidden="true"
              />
            ) : null}
          </div>

          <div className="mt-6 flex shrink-0 items-center justify-between border-t pt-4">
            {!isFirstStep ? (
              <Button
                type="button"
                onClick={goBack}
                disabled={isBusy}
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
                disabled={isBusy}
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
              disabled={isBusy || (isLastStep && !canSubmit)}
              title={
                isLastStep && !canSubmit
                  ? "Nenhuma alteração para salvar"
                  : undefined
              }
              aria-busy={isBusy}
              className="p-4 text-sm font-semibold"
            >
              {isBusy ? (
                <>
                  <Loader2 className="size-4 animate-spin" />
                  <span>{isLastStep ? "Salvando..." : "Processando..."}</span>
                </>
              ) : (
                <>
                  <span>{isLastStep ? submitLabel : "Avançar"}</span>
                  {!isLastStep ? <ArrowRight className="size-4" /> : null}
                </>
              )}
            </Button>
          </div>
        </form>
      </div>
    </ImageUploadStateProvider>
  );
}
