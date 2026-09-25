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
  return (
    <div
      class={cn(
        "w-full max-w-full min-w-0 rounded-xl border border-border/70 bg-card p-3.5 sm:p-4 space-y-3.5 sm:space-y-4 shadow-xs overflow-hidden",
        props.class,
      )}
    >
      <Show when={props.summary.title || props.summary.badge}>
        <div class="flex items-center justify-between gap-2 border-b border-border/50 pb-2.5">
          <Show when={props.summary.title}>
            <h3 class="text-xs sm:text-sm font-semibold text-foreground truncate">
              {props.summary.title}
            </h3>
          </Show>
          <Show when={props.summary.badge}>
            <span class="shrink-0 rounded-md bg-primary/10 px-2 py-0.5 text-[10px] sm:text-[11px] font-medium text-primary">
              {props.summary.badge}
            </span>
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
              typeof item.href === "function" ? (item.href as () => string | undefined)() : item.href;

            return (
              <div class={cn("min-w-0 w-full overflow-hidden", item.fullWidth && "sm:col-span-2")}>
                <span class="block text-[11px] font-medium text-muted-foreground truncate">
                  {item.label}
                </span>
                <Show
                  when={hrefVal()}
                  fallback={
                    <span
                      class={cn(
                        "mt-0.5 block text-foreground truncate max-w-full",
                        item.mono && "font-mono text-[11px] break-all",
                        item.fullWidth &&
                          "whitespace-pre-wrap break-words [overflow-wrap:anywhere] leading-relaxed text-foreground/90 font-normal",
                      )}
                    >
                      {displayVal()}
                    </span>
                  }
                >
                  <a
                    href={hrefVal()!}
                    target="_blank"
                    rel="noopener noreferrer"
                    class={cn(
                      "mt-0.5 block text-primary hover:underline truncate max-w-full",
                      item.mono && "font-mono text-[11px] break-all",
                    )}
                  >
                    {displayVal()}
                  </a>
                </Show>
              </div>
            );
          }}
        </For>

        <Show when={props.summary.extra}>
          <div class="col-span-1 sm:col-span-2 pt-1 w-full max-w-full min-w-0">
            {props.summary.extra?.()}
          </div>
        </Show>
      </div>
    </div>
  );
}
