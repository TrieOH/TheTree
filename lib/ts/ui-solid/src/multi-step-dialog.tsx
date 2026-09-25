import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal } from "solid-js";
import { Button } from "./button";
import { Dialog, type DialogSize } from "./dialog";
import { cn } from "./lib/cn";
import {
  MultiStepFields,
  MultiStepStepper,
  MultiStepSummaryCard,
  type MultiStepItem,
  type MultiStepRenderContext,
} from "./multi-step";

export * from "./multi-step";

export interface MultiStepDialogProps<T = Record<string, unknown>> {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: JSX.Element;
  description?: JSX.Element;
  steps: MultiStepItem<T>[];
  currentStep?: number;
  onStepChange?: (step: number) => void;
  values?: T;
  errors?: Partial<Record<keyof T & string, string>>;
  onChange?: (key: keyof T & string, value: unknown) => void;
  children?: (context: MultiStepRenderContext<T>) => JSX.Element;
  /** Async or sync callback triggered when submitting the last step. Returning false keeps modal open. */
  onSubmit?: () => Promise<boolean | void> | boolean | void;
  /** Optional custom form submit handler (e.g. for requestSubmit or Enter key) */
  onFormSubmit?: (event: SubmitEvent) => Promise<boolean | void> | boolean | void;
  /** Optional guard before advancing. Return false to block moving forward. */
  onBeforeNext?: (currentStep: number) => Promise<boolean> | boolean;
  loading?: boolean;
  submitLabel?: string;
  nextLabel?: string;
  backLabel?: string;
  cancelLabel?: string;
  onCancel?: () => void;
  formId?: string;
  closable?: boolean;
  size?: DialogSize;
  /**
   * When true (default for MultiStepDialog), locks dialog panel height to `calc(100dvh - 2rem)`
   * to eliminate jarring layout shifts between steps with different content heights.
   */
  fixedHeight?: boolean;
  class?: string;
  contentClass?: string;
  preventClose?: boolean;
}

export function MultiStepDialog<T = Record<string, unknown>>(
  props: MultiStepDialogProps<T>,
): JSX.Element {
  const [internalStep, setInternalStep] = createSignal(0);

  // Sync internal step if controlled externally
  createEffect(
    () => props.currentStep,
    (currentStep) => {
      if (currentStep !== undefined) {
        setInternalStep(currentStep);
      }
    },
  );

  // Reset to step 0 when modal closes if not controlled (both dependencies tracked in compute)
  createEffect(
    () => [props.open, props.currentStep] as const,
    ([open, currentStep]) => {
      if (!open && currentStep === undefined) {
        setInternalStep(0);
      }
    },
  );

  const stepIndex = () =>
    props.currentStep !== undefined ? props.currentStep : internalStep();

  const totalSteps = () => props.steps.length;
  const isFirst = () => stepIndex() <= 0;
  const isLast = () => stepIndex() >= totalSteps() - 1;
  const currentStepItem = () => props.steps[stepIndex()] ?? props.steps[0];

  const setStep = (next: number) => {
    const clamped = Math.max(0, Math.min(next, totalSteps() - 1));
    if (props.onStepChange) {
      props.onStepChange(clamped);
    } else {
      setInternalStep(clamped);
    }
  };

  const handleNext = async () => {
    if (props.loading) return;

    if (props.onBeforeNext) {
      const allowed = await props.onBeforeNext(stepIndex());
      if (!allowed) return;
    }

    if (!isLast()) {
      setStep(stepIndex() + 1);
    } else if (props.onSubmit) {
      const result = await props.onSubmit();
      if (result !== false) {
        props.onOpenChange(false);
      }
    }
  };

  const handleBack = () => {
    if (props.loading) return;
    if (!isFirst()) {
      setStep(stepIndex() - 1);
    } else if (props.onCancel) {
      props.onCancel();
    } else {
      props.onOpenChange(false);
    }
  };

  const handleFormSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    if (props.onFormSubmit) {
      const result = await props.onFormSubmit(event);
      if (result !== false && isLast()) {
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
      preventClose={props.preventClose ?? props.loading}
      header={
        <div class="flex flex-col gap-3.5 border-b border-border/40 px-4 py-3.5 sm:px-6 sm:py-4 bg-card shrink-0 min-w-0">
          <div class="flex items-start justify-between gap-3 min-w-0">
            <div class="flex min-w-0 flex-1 flex-col pr-2">
              <h2 class="text-base sm:text-lg font-semibold tracking-tight text-foreground leading-snug truncate">
                {props.title}
              </h2>
              <Show when={props.description}>
                {(desc) => (
                  <p class="text-xs sm:text-sm text-muted-foreground mt-0.5 leading-relaxed">
                    {desc()}
                  </p>
                )}
              </Show>
            </div>

            <Show when={props.closable !== false}>
              <button
                type="button"
                aria-label="Fechar"
                class="-mr-1.5 -mt-1 inline-flex size-8 shrink-0 items-center justify-center rounded-xl text-muted-foreground/70 transition-colors hover:bg-muted hover:text-foreground active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring cursor-pointer"
                onClick={() => props.onOpenChange(false)}
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  class="size-4"
                  aria-hidden="true"
                >
                  <path d="M18 6 6 18M6 6l12 12" stroke-linecap="round" />
                </svg>
              </button>
            </Show>
          </div>

          <MultiStepStepper
            steps={props.steps}
            currentStep={stepIndex()}
            onStepClick={(target) => {
              if (target < stepIndex()) setStep(target);
            }}
          />
        </div>
      }
      footer={
        <div class="flex w-full items-center justify-between gap-2.5 min-w-0">
          <Button
            type="button"
            variant="outline"
            onClick={handleBack}
            disabled={props.loading}
            class="h-10 sm:h-9 px-3 text-xs sm:text-sm"
          >
            <Show when={!isFirst()}>
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                class="size-3.5 mr-1"
                aria-hidden="true"
              >
                <path d="M19 12H5M12 19l-7-7 7-7" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </Show>
            <span>{backLabel()}</span>
          </Button>

          <div class="flex items-center gap-2">
            <Show
              when={isLast()}
              fallback={
                <Button
                  type="button"
                  variant="default"
                  onClick={() => void handleNext()}
                  disabled={props.loading}
                  class="h-10 sm:h-9 min-w-24 px-4 text-xs sm:text-sm"
                >
                  <span class="inline-flex items-center gap-1.5">
                    <span>{nextLabel()}</span>
                    <svg
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      class="size-3.5"
                      aria-hidden="true"
                    >
                      <path d="M5 12h14M12 5l7 7-7 7" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                  </span>
                </Button>
              }
            >
              <Button
                type="button"
                variant="default"
                onClick={() => void handleNext()}
                disabled={props.loading}
                class="h-10 sm:h-9 min-w-24 px-4 text-xs sm:text-sm"
              >
                <Show
                  when={props.loading}
                  fallback={<span>{submitLabel()}</span>}
                >
                  <span class="inline-flex items-center gap-1.5">
                    <svg
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.5"
                      class="size-3.5 animate-spin"
                      aria-hidden="true"
                    >
                      <circle cx="12" cy="12" r="10" stroke-opacity="0.25" />
                      <path d="M12 2a10 10 0 0 1 10 10" stroke-linecap="round" />
                    </svg>
                    <span>Salvando...</span>
                  </span>
                </Show>
              </Button>
            </Show>
          </div>
        </div>
      }
    >
      <div class="min-h-0 flex-1 w-full max-w-full min-w-0">
        <Show
          when={hasCustomChildren()}
          fallback={
            <form
              id={props.formId}
              class="flex flex-col gap-4 w-full max-w-full min-w-0"
              onSubmit={handleFormSubmit}
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
          }
        >
          {props.children?.(renderContext)}
        </Show>
      </div>
    </Dialog>
  );
}
