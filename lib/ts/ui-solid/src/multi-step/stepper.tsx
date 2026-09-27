import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { cn } from "../lib/cn";
import type { MultiStepStepperProps } from "./types";

/**
 * Responsive Stepper:
 * - Mobile (< 640px): Compact progress indicator with progress bar & step title
 * - Desktop (≥ 640px): Step pills with numbers/checkmarks, titles & connectors
 */
export function MultiStepStepper(props: MultiStepStepperProps): JSX.Element {
  const total = () => props.steps.length;
  const safeCurrent = () => Math.min(Math.max(0, props.currentStep), total() - 1);
  const activeStep = () => props.steps[safeCurrent()];
  const progressPercent = () =>
    total() > 1 ? Math.round((safeCurrent() / (total() - 1)) * 100) : 100;

  return (
    <div class={cn("w-full space-y-2.5 min-w-0", props.class)}>
      {/* Mobile view (< 640px) */}
      <div class="space-y-1.5 sm:hidden min-w-0">
        <div class="flex items-center justify-between gap-2 text-xs">
          <span class="font-medium text-foreground truncate min-w-0">
            {activeStep()?.title}
          </span>
          <span class="shrink-0 font-mono text-[11px] text-muted-foreground">
            {safeCurrent() + 1} de {total()}
          </span>
        </div>
        <div
          class="h-1.5 w-full overflow-hidden rounded-full bg-muted/60"
          role="progressbar"
          aria-valuenow={safeCurrent() + 1}
          aria-valuemin={1}
          aria-valuemax={total()}
          aria-label={`Passo ${safeCurrent() + 1} de ${total()}: ${activeStep()?.title ?? ""}`}
        >
          <div
            class="h-full bg-primary transition-all duration-300 ease-out"
            style={{ width: `${progressPercent()}%` }}
          />
        </div>
      </div>

      {/* Desktop view (≥ 640px) */}
      <ol class="hidden sm:flex sm:items-center sm:gap-2 text-xs min-w-0">
        <For each={props.steps}>
          {(step, idx) => {
            const isCompleted = () => idx() < safeCurrent();
            const isActive = () => idx() === safeCurrent();
            const isClickable = () =>
              Boolean(props.onStepClick && idx() <= safeCurrent());

            return (
              <li class="flex flex-1 items-center gap-2 min-w-0">
                <button
                  type="button"
                  disabled={!isClickable()}
                  onClick={() => props.onStepClick?.(idx())}
                  class={cn(
                    "group flex items-center gap-2 rounded-md text-left transition-colors min-w-0",
                    isClickable() ? "cursor-pointer hover:opacity-85" : "cursor-default",
                  )}
                  aria-current={isActive() ? "step" : undefined}
                >
                  <span
                    class={cn(
                      "flex size-6 shrink-0 items-center justify-center rounded-full text-[11px] font-semibold transition-all",
                      isCompleted()
                        ? "bg-primary text-primary-foreground shadow-xs"
                        : isActive()
                          ? "border-2 border-primary bg-primary/10 text-primary font-bold shadow-xs"
                          : "border border-border bg-muted/50 text-muted-foreground",
                    )}
                  >
                    <Show
                      when={isCompleted()}
                      fallback={<span>{idx() + 1}</span>}
                    >
                      <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="3"
                        class="size-3"
                        aria-hidden="true"
                      >
                        <polyline points="20 6 9 17 4 12" stroke-linecap="round" stroke-linejoin="round" />
                      </svg>
                    </Show>
                  </span>

                  <span class="min-w-0">
                    <span
                      class={cn(
                        "block truncate font-medium text-xs",
                        isActive()
                          ? "text-foreground font-semibold"
                          : isCompleted()
                            ? "text-foreground"
                            : "text-muted-foreground",
                      )}
                    >
                      {step.title}
                    </span>
                    <Show when={step.description}>
                      {(desc) => (
                        <span class="block truncate text-[10px] text-muted-foreground">
                          {desc()}
                        </span>
                      )}
                    </Show>
                  </span>
                </button>

                {/* Connector line between steps */}
                <Show when={idx() < total() - 1}>
                  <div
                    class={cn(
                      "h-px flex-1 transition-colors min-w-2",
                      idx() < safeCurrent() ? "bg-primary/50" : "bg-border/60",
                    )}
                    aria-hidden="true"
                  />
                </Show>
              </li>
            );
          }}
        </For>
      </ol>
    </div>
  );
}
