import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { cn } from "../lib/cn";
import type { MultiStepSummaryConfig } from "./types";

/**
 * Modular subcomponent for reviewing and confirming configured values before submission.
 * Completely responsive and avoids any mobile layout overflow.
 */
export function MultiStepSummaryCard(props: {
  summary: MultiStepSummaryConfig;
  class?: string;
}): JSX.Element {
  const badgeText = () =>
    typeof props.summary.badge === "function"
      ? props.summary.badge()
      : props.summary.badge;

  return (
    <div
      class={cn(
        "w-full max-w-full min-w-0 rounded-xl border border-border/70 bg-card p-3.5 sm:p-4 space-y-3.5 sm:space-y-4 shadow-xs overflow-hidden",
        props.class,
      )}
    >
      <Show when={props.summary.title || badgeText()}>
        <div class="flex items-center justify-between gap-2 border-b border-border/50 pb-2.5">
          <Show when={props.summary.title}>
            <h3 class="text-xs sm:text-sm font-semibold text-foreground truncate">
              {props.summary.title}
            </h3>
          </Show>
          <Show when={badgeText()}>
            {(text) => (
              <span class="shrink-0 rounded-md bg-primary/10 px-2 py-0.5 text-[10px] sm:text-[11px] font-medium text-primary">
                {text()}
              </span>
            )}
          </Show>
        </div>
      </Show>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-3.5 text-xs w-full max-w-full min-w-0">
        <For each={props.summary.items}>
          {(item) => {
            const rawVal = () =>
              typeof item.value === "function" ? (item.value as () => unknown)() : item.value;
            const displayVal = () => {
              const v = rawVal();
              return v != null && v !== "" ? String(v) : "—";
            };
            const hrefVal = () =>
              typeof item.href === "function" ? item.href() : item.href;

            return (
              <div
                class={cn(
                  "min-w-0 flex flex-col gap-1 p-2 rounded-lg bg-muted/40 border border-border/30",
                  item.fullWidth && "sm:col-span-2",
                )}
              >
                <span class="text-[11px] font-medium text-muted-foreground truncate">
                  {item.label}
                </span>
                <Show
                  when={hrefVal()}
                  fallback={
                    <span
                      class={cn(
                        "text-xs text-foreground truncate font-medium",
                        item.mono && "font-mono text-[11px]",
                      )}
                    >
                      {displayVal()}
                    </span>
                  }
                >
                  {(href) => (
                    <a
                      href={href()}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-xs text-primary underline-offset-4 hover:underline truncate font-mono text-[11px]"
                    >
                      {displayVal()}
                    </a>
                  )}
                </Show>
              </div>
            );
          }}
        </For>
      </div>

      <Show when={props.summary.extra}>
        {(renderExtra) => <div class="w-full min-w-0 pt-1">{renderExtra()()}</div>}
      </Show>
    </div>
  );
}
