import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { cn } from "../lib/cn";
import { shapeClass } from "../lib/shape";
import type { MultiStepStepperProps } from "./types";

/** Compact progress with optional navigation to previously completed steps. */
export function MultiStepStepper(props: MultiStepStepperProps): JSX.Element {
  const total = () => props.steps.length;
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
          aria-valuetext={`Etapa ${current() + 1} de ${total()}: ${props.steps[current()]?.title ?? ""}`}
        />
        <div class={cn("hidden", automatic() && "@max-[28rem]/stepper:block")}>
          <div class="flex min-w-0 items-center justify-between gap-3 text-xs">
            <span class="truncate font-medium text-foreground">
              {props.steps[current()]?.title}
            </span>
            <span class="shrink-0 tabular-nums text-muted-foreground">
              {current() + 1} de {total()}
            </span>
          </div>
          <ol class="mt-1 flex gap-1.5" aria-label="Progresso das etapas">
            <For each={props.steps}>
              {(step, index) => (
                <li class="min-w-0 flex-1">
                  <button
                    type="button"
                    disabled={!props.onStepClick || index() >= current()}
                    onClick={() => {
                      if (index() < current()) props.onStepClick?.(index());
                    }}
                    aria-label={`${index() + 1}. ${step.title}. ${index() < current() ? "Concluída" : index() === current() ? "Etapa atual" : "Pendente"}`}
                    aria-current={index() === current() ? "step" : undefined}
                    class="flex min-h-10 w-full items-center rounded-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-default"
                  >
                    <span
                      aria-hidden="true"
                      class={cn(
                        "h-1 w-full rounded-full",
                        index() <= current() ? "bg-accent" : "bg-border",
                      )}
                    />
                  </button>
                </li>
              )}
            </For>
          </ol>
          <Show
            when={props.showDescriptions && props.steps[current()]?.description}
          >
            <p class="text-xs leading-5 text-muted-foreground">
              {props.steps[current()]?.description}
            </p>
          </Show>
        </div>
        <ol class={cn(
          "m-0 list-none p-0",
          vertical() ? "grid gap-0" : "flex items-center gap-3",
          automatic() && "@max-[28rem]/stepper:hidden",
        )}>
          <For each={props.steps}>
            {(step, index) => {
              const completed = () => index() < current();
              const active = () => index() === current();
              const clickable = () => !!props.onStepClick && completed();
              const status = () =>
                completed()
                  ? "Concluída"
                  : active()
                    ? "Etapa atual"
                    : "Pendente";

              return (
                <li class={cn(
                  "min-w-0",
                  vertical()
                    ? "relative block pb-4 last:pb-0"
                    : "flex flex-1 items-center gap-3 last:max-w-[45%] last:flex-[0_1_auto]",
                )}>
                  <button
                    type="button"
                    disabled={!clickable()}
                    onClick={() => {
                      if (clickable()) props.onStepClick?.(index());
                    }}
                    aria-current={active() ? "step" : undefined}
                    aria-label={`${index() + 1}. ${step.title}. ${status()}`}
                    class={cn(
                      "group flex min-h-10 w-auto min-w-0 items-center gap-2 rounded-md text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
                      clickable() ? "cursor-pointer" : "cursor-default",
                    )}
                  >
                    <span
                      aria-hidden="true"
                      class={cn(
                        "flex size-6 shrink-0 items-center justify-center border text-xs font-semibold tabular-nums transition-colors motion-reduce:transition-none",
                        shapeClass(props.shape ?? "pill"),
                        completed()
                          ? "border-accent/30 bg-accent/15 text-foreground group-hover:border-accent"
                          : active()
                            ? "border-accent bg-accent text-accent-foreground"
                            : "border-border bg-muted text-muted-foreground",
                      )}
                    >
                      <Show when={completed()} fallback={index() + 1}>
                        <svg
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          class="size-3.5"
                        >
                          <path
                            d="m5 12 4 4L19 6"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                          />
                        </svg>
                      </Show>
                    </span>
                    <span class="min-w-0">
                      <span
                        class={cn(
                          "block truncate text-xs leading-5",
                          active()
                            ? "font-semibold text-foreground"
                            : completed()
                              ? "font-medium text-foreground"
                              : "font-medium text-muted-foreground",
                        )}
                      >
                        {step.title}
                      </span>
                      <Show when={props.showDescriptions && step.description}>
                        <span class="mt-0.5 block truncate text-[11px] leading-4 text-muted-foreground">
                          {step.description}
                        </span>
                      </Show>
                    </span>
                  </button>
                  <Show when={index() < total() - 1}>
                    <span
                      aria-hidden="true"
                      class={cn(
                        "pointer-events-none",
                        vertical()
                          ? "absolute top-9 bottom-1 left-[calc(0.75rem-0.5px)] h-auto w-px min-w-0"
                          : "h-px min-w-4 flex-1",
                        completed() ? "bg-accent/60" : "bg-border",
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
