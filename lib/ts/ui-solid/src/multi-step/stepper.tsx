import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { cn } from "../lib/cn";
import { shapeClass } from "../lib/shape";
import type { MultiStepStepperProps } from "./types";

/** Compact progress with optional navigation to previously completed steps. */
export function MultiStepStepper(props: MultiStepStepperProps): JSX.Element {
  const total = () => props.steps?.length ?? 0;
  const current = () => Math.max(0, Math.min(props.currentStep, total() - 1));
  const vertical = () => props.orientation === "vertical";
  const automatic = () => !props.orientation || props.orientation === "auto";

  return (
    <Show when={total() > 0}>
      <nav
        aria-label={props.ariaLabel ?? "Etapas"}
        data-orientation={props.orientation ?? "auto"}
        class={cn("@container/stepper w-full min-w-0", props.class)}
      >
        <div
          role="progressbar"
          class="sr-only"
          aria-label={props.ariaLabel ?? "Progresso das etapas"}
          aria-valuemin={1}
          aria-valuemax={total()}
          aria-valuenow={current() + 1}
          aria-valuetext={`Etapa ${current() + 1} de ${total()}: ${props.steps?.[current()]?.title ?? ""}`}
        />
        <div class={cn("hidden", automatic() && "@max-[28rem]/stepper:block")}>
          <div class="flex min-w-0 items-center justify-between gap-3 text-xs">
            <span class="truncate font-medium text-foreground">
              {props.steps?.[current()]?.title}
            </span>
            <span class="shrink-0 tabular-nums text-muted-foreground">
              {current() + 1} de {total()}
            </span>
          </div>
          <ol class="mt-1 flex gap-1.5" aria-label="Progresso das etapas">
            <For each={props.steps ?? []}>
              {(_, i) => (
                <li class="h-1 flex-1 rounded-full overflow-hidden bg-muted">
                  <div
                    class={cn(
                      "h-full transition-all",
                      i() <= current() ? "bg-primary" : "bg-transparent",
                    )}
                  />
                </li>
              )}
            </For>
          </ol>
        </div>

        <ol
          class={cn(
            vertical() ? "grid gap-3" : "flex items-center gap-3",
            automatic() && "@max-[28rem]/stepper:hidden",
          )}
        >
          <For each={props.steps ?? []}>
            {(step, i) => {
              const isDone = () => i() < current();
              const isCurrent = () => i() === current();
              const isInteractive = () =>
                typeof props.onStepClick === "function" && isDone();

              return (
                <li
                  class={cn(
                    "flex flex-1 items-center gap-3",
                    vertical() && "w-full",
                  )}
                >
                  <button
                    type="button"
                    disabled={!isInteractive()}
                    onClick={() => props.onStepClick?.(i())}
                    class={cn(
                      "group flex items-center gap-3 text-left transition-colors",
                      isInteractive() && "cursor-pointer",
                    )}
                  >
                    <span
                      class={cn(
                        "flex size-7 shrink-0 items-center justify-center text-xs font-semibold border transition-all",
                        shapeClass(props.shape ?? "pill"),
                        isCurrent() &&
                          "border-primary bg-primary text-primary-foreground shadow-xs ring-2 ring-primary/20",
                        isDone() &&
                          "border-primary/40 bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground",
                        !isCurrent() &&
                          !isDone() &&
                          "border-border bg-muted/40 text-muted-foreground",
                      )}
                    >
                      <Show when={isDone()} fallback={i() + 1}>
                        <svg
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2.5"
                          class="size-3.5"
                          aria-hidden="true"
                        >
                          <path
                            d="M20 6 9 17l-5-5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                          />
                        </svg>
                      </Show>
                    </span>

                    <span class="flex min-w-0 flex-col">
                      <span
                        class={cn(
                          "truncate text-xs font-medium",
                          isCurrent() && "text-foreground font-semibold",
                          isDone() &&
                            "text-muted-foreground group-hover:text-foreground",
                          !isCurrent() &&
                            !isDone() &&
                            "text-muted-foreground/80",
                        )}
                      >
                        {step.title}
                      </span>
                      <Show when={props.showDescriptions && step.description}>
                        {(desc) => (
                          <span class="truncate text-[11px] text-muted-foreground/70">
                            {desc()}
                          </span>
                        )}
                      </Show>
                    </span>
                  </button>

                  <Show when={!vertical() && i() < total() - 1}>
                    <div
                      aria-hidden="true"
                      class={cn(
                        "h-px flex-1 transition-colors",
                        isDone() ? "bg-primary/50" : "bg-border/60",
                      )}
                    />
                  </Show>
                </li>
              );
            }}
          </For>
        </ol>
      </nav>
    </Show>
  );
}
