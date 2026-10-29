import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal } from "solid-js";
import { Button } from "./button";
import { Dialog, type DialogProps } from "./dialog";
import { cn } from "./lib/cn";
import { MultiStepFields } from "./multi-step/fields";
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
  MultiStepStepperStep,
  MultiStepSummaryConfig,
  MultiStepSummaryItem,
  TextMultiStepField,
  TextareaMultiStepField,
} from "./multi-step/types";

export function MultiStepDialog<T = Record<string, unknown>>(
  props: MultiStepDialogProps<T> & {
    items?: MultiStepItem<T>[];
    size?: DialogProps["size"];
    backLabel?: string;
  },
): JSX.Element {
  const [internalStep, setInternalStep] = createSignal(0);
  const [submissionPending, setSubmissionPending] = createSignal(false);
  // Signals update the UI in a batch; this lock also covers same-turn submits.
  let submissionInFlight = false;
  const busy = () => Boolean(props.loading || submissionPending() || submissionInFlight);

  const runSubmission = async (action: () => Promise<void>) => {
    if (busy()) return;
    submissionInFlight = true;
    setSubmissionPending(true);
    try {
      await action();
    } finally {
      submissionInFlight = false;
      setSubmissionPending(false);
    }
  };

  const stepItems = () => props.steps ?? props.items ?? [];
  const totalSteps = () => stepItems().length;

  const isControlled = () => props.currentStep !== undefined;
  const stepIndex = () => props.currentStep ?? internalStep();

  const setStep = (next: number) => {
    const clamped = Math.max(0, Math.min(next, totalSteps() - 1));
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

  const isFirst = () => stepIndex() === 0;
  const isLast = () => stepIndex() === totalSteps() - 1;

  const currentStepItem = () => stepItems()[stepIndex()] ?? stepItems()[0];

  const handleNext = async () => {
    if (busy()) return;
    if (props.onBeforeNext) {
      const allowed = await props.onBeforeNext(stepIndex());
      if (allowed === false) return;
    }

    if (isLast()) {
      const submitFn = props.onFormSubmit ?? props.onSubmit;
      await runSubmission(async () => {
        const result = await submitFn?.();
        if (result !== false) {
          props.onOpenChange(false);
        }
      });
      return;
    }

    setStep(stepIndex() + 1);
  };

  const handleBack = () => {
    if (busy()) return;
    if (isFirst()) {
      props.onOpenChange(false);
      return;
    }
    setStep(stepIndex() - 1);
  };

  const handleFormSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    const submitFn = props.onFormSubmit;
    if (submitFn) {
      await runSubmission(async () => {
        const result = await submitFn();
        if (result !== false) {
          props.onOpenChange(false);
        }
      });
      return;
    }
    await handleNext();
  };

  const submitLabel = () => props.submitLabel ?? "Salvar";
  const nextLabel = () => props.nextLabel ?? "Continuar";
  const backLabel = () =>
    props.backLabel ??
    (isFirst() ? (props.cancelLabel ?? "Cancelar") : "Voltar");

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
    values: () => (props.values ?? ({} as T)),
    errors: () =>
      (props.errors ?? {}) as Partial<Record<keyof T & string, string>>,
    next: () => void handleNext(),
    prev: handleBack,
    cancel: () => props.onOpenChange(false),
    onChange: (key, val) => props.onChange?.(key, val),
    goNext: () => void handleNext(),
    goBack: handleBack,
    goTo: setStep,
  };

  const hasCustomChildren = () => typeof props.children === "function";
  const activeSummary = () => currentStepItem()?.summary;
  const activeFields = () => currentStepItem()?.fields;

  return (
    <Dialog
      shape={props.shape}
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
                disabled={busy()}
              >
                {backLabel()}
              </Button>

              <div class="flex items-center gap-2">
                <Button
                  type="button"
                  variant="default"
                  onClick={handleNext}
                  disabled={busy()}
                  class="min-w-28 font-medium transition-all shadow-xs"
                >
                  {isLast() ? submitLabel() : nextLabel()}
                </Button>
              </div>
            </div>
          }
        >
          {typeof props.footer === "function"
            ? props.footer(renderContext)
            : props.footer}
        </Show>
      }
    >
      <div class="flex flex-col gap-5">
        {/* Stepper progress (desktop and mobile responsive) */}
        <MultiStepStepper
          steps={stepItems()}
          currentStep={stepIndex()}
          onStepClick={(target) => {
            if (!busy() && target < stepIndex()) {
              setStep(target);
            }
          }}
        />

        {/* Step Body */}
        <div class="min-h-0">
          <Show
            when={!hasCustomChildren()}
            fallback={
              typeof props.children === "function"
                ? props.children(renderContext)
                : props.children
            }
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
                    fallback={currentStepItem()?.render?.(renderContext) ?? null}
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
                {(summary) => <MultiStepSummaryCard summary={summary()} />}
              </Show>
            </form>
          </Show>
        </div>
      </div>
    </Dialog>
  );
}
