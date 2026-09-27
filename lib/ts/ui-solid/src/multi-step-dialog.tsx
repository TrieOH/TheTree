import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal } from "solid-js";
import { Button } from "./button";
import { Dialog, type DialogProps } from "./dialog";
import { cn } from "./lib/cn";
import {
  MultiStepCustomField,
  MultiStepFields,
  MultiStepTextField,
  MultiStepTextareaField,
} from "./multi-step/fields";
import { MultiStepStepper } from "./multi-step/stepper";
import { MultiStepSummaryCard } from "./multi-step/summary-card";
import type {
  MultiStepDialogProps,
  MultiStepItem,
  MultiStepRenderContext,
} from "./multi-step/types";

export {
  MultiStepFields,
  MultiStepTextField,
  MultiStepTextareaField,
  MultiStepCustomField,
} from "./multi-step/fields";
export { MultiStepStepper } from "./multi-step/stepper";
export { MultiStepSummaryCard } from "./multi-step/summary-card";
export type {
  BaseMultiStepField,
  CustomMultiStepField,
  MultiStepDialogProps,
  MultiStepField,
  MultiStepFieldConfig,
  MultiStepFieldKind,
  MultiStepFieldsProps,
  MultiStepItem,
  MultiStepRenderContext,
  MultiStepStepperProps,
  MultiStepSummaryConfig,
  MultiStepSummaryItem,
  TextMultiStepField,
  TextareaMultiStepField,
} from "./multi-step/types";

export function MultiStepDialog<T = Record<string, unknown>>(
  props: MultiStepDialogProps<T> & {
    size?: DialogProps["size"];
    backLabel?: string;
  },
): JSX.Element {
  const [internalStep, setInternalStep] = createSignal(0);

  const isControlled = () => props.currentStep !== undefined;
  const stepIndex = () => (isControlled() ? props.currentStep! : internalStep());

  const setStep = (next: number) => {
    const clamped = Math.max(0, Math.min(next, props.steps.length - 1));
    if (!isControlled()) {
      setInternalStep(clamped);
    }
    props.onStepChange?.(clamped);
  };

  createEffect(
    () => [props.open, props.currentStep] as const,
    ([open, currentStep]) => {
      if (!open && currentStep === undefined) {
        setInternalStep(0);
      } else if (currentStep !== undefined) {
        setInternalStep(currentStep);
      }
    },
  );

  const totalSteps = () => props.steps.length;
  const isFirst = () => stepIndex() === 0;
  const isLast = () => stepIndex() === totalSteps() - 1;

  const currentStepItem = () => props.steps[stepIndex()] ?? props.steps[0];

  const handleNext = async () => {
    if (props.loading) return;

    if (props.onBeforeNext) {
      const allowed = await props.onBeforeNext(stepIndex());
      if (allowed === false) return;
    }

    if (isLast()) {
      const submitFn = props.onFormSubmit ?? props.onSubmit;
      const result = await submitFn?.();
      if (result !== false) {
        props.onOpenChange(false);
      }
      return;
    }

    setStep(stepIndex() + 1);
  };

  const handleBack = () => {
    if (props.loading) return;
    if (isFirst()) {
      props.onOpenChange(false);
      return;
    }
    setStep(stepIndex() - 1);
  };

  const handleFormSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    if (props.onFormSubmit) {
      const result = await props.onFormSubmit();
      if (result !== false) {
        props.onOpenChange(false);
      }
      return;
    }
    await handleNext();
  };

  const submitLabel = () => props.submitLabel ?? "Salvar";
  const nextLabel = () => props.nextLabel ?? "Continuar";
  const backLabel = () =>
    props.backLabel ?? (isFirst() ? (props.cancelLabel ?? "Cancelar") : "Voltar");

  const renderContext: MultiStepRenderContext<T> = {
    get currentStep() {
      return stepIndex();
    },
    get step() {
      return currentStepItem();
    },
    get isFirst() {
      return isFirst();
    },
    get isLast() {
      return isLast();
    },
    values: () => (props.values ?? {}) as T,
    errors: () => (props.errors ?? {}) as Partial<Record<keyof T & string, string>>,
    next: () => void handleNext(),
    prev: handleBack,
    cancel: () => props.onOpenChange(false),
    onChange: (key, val) => props.onChange?.(key, val),
    goNext: () => void handleNext(),
    goBack: handleBack,
    goTo: setStep,
  };

  const hasCustomChildren = () => typeof props.children === "function";
  const activeSummary = () => currentStepItem().summary;
  const activeFields = () => currentStepItem().fields;

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      size={props.size ?? "lg"}
      fixedHeight={props.fixedHeight ?? true}
      class={props.class}
      contentClass={cn("p-4 sm:p-6", props.contentClass)}
      title={props.title}
      description={props.description}
      footer={
        <Show
          when={props.footer}
          fallback={
            <div class="flex items-center justify-between gap-3 pt-2">
              <Button
                type="button"
                variant="ghost"
                onClick={handleBack}
                disabled={props.loading}
              >
                {backLabel()}
              </Button>

              <div class="flex items-center gap-2">
                <Button
                  type="button"
                  variant="default"
                  onClick={handleNext}
                  disabled={props.loading}
                  class="min-w-28 font-medium transition-all shadow-xs"
                >
                  {isLast() ? submitLabel() : nextLabel()}
                </Button>
              </div>
            </div>
          }
        >
          {props.footer!(renderContext)}
        </Show>
      }
    >
      <div class="flex flex-col gap-5">
        {/* Stepper progress (desktop and mobile responsive) */}
        <MultiStepStepper
          steps={props.steps as MultiStepItem<any>[]}
          currentStep={stepIndex()}
          onStepClick={(target) => {
            if (target < stepIndex()) {
              setStep(target);
            }
          }}
        />

        {/* Step Body */}
        <div class="min-h-0">
          <Show
            when={!hasCustomChildren()}
            fallback={props.children!(renderContext)}
          >
            <form
              id={props.formId ?? "multi-step-form"}
              onSubmit={handleFormSubmit}
              class="space-y-4"
            >
              <Show
                when={activeSummary()}
                fallback={
                  <Show
                    when={activeFields()}
                    fallback={currentStepItem().render?.(renderContext) ?? null}
                  >
                    <MultiStepFields
                      fields={activeFields() ?? []}
                      values={props.values}
                      errors={props.errors}
                      onChange={props.onChange}
                    />
                  </Show>
                }
              >
                <MultiStepSummaryCard summary={activeSummary()!} />
              </Show>
            </form>
          </Show>
        </div>
      </div>
    </Dialog>
  );
}
