import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import CircleCheck from "~icons/lucide/circle-check";
import { Loading } from "./loading";
import { cn } from "./lib/cn";
export interface TimelineItem {
  id: string;
  title: string;
  date: Date | string;
  description?: JSX.Element;
  icon?: JSX.Element;
  tone?: "default" | "success" | "warning" | "destructive";
  action?: JSX.Element;
}
export interface TimelineProps {
  items: readonly TimelineItem[];
  orientation?: "vertical" | "horizontal";
  loading?: boolean;
  emptyMessage?: string;
  dateFormatter?: (date: Date) => string;
  class?: string;
  ariaLabel?: string;
}
const tones = {
  default: "border-border bg-muted text-muted-foreground",
  success:
    "border-emerald-500/20 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
  warning:
    "border-amber-500/20 bg-amber-500/10 text-amber-600 dark:text-amber-400",
  destructive: "border-destructive/20 bg-destructive/10 text-destructive",
};
export function Timeline(props: TimelineProps) {
  const horizontal = () => props.orientation === "horizontal";
  const date = (item: TimelineItem) => {
    const value = new Date(item.date);
    return Number.isFinite(value.getTime()) ? value : null;
  };
  return (
    <div
      class={cn("w-full min-w-0", props.class)}
      aria-busy={props.loading ? "true" : "false"}
    >
      <Show when={!props.loading} fallback={<Loading size="sm" tone="muted" />}>
        <Show
          when={props.items.length}
          fallback={
            <p class="py-6 text-center text-sm text-muted-foreground">
              {props.emptyMessage ?? "Nenhuma atividade registrada."}
            </p>
          }
        >
          <ol
            aria-label={props.ariaLabel ?? "Histórico de atividades"}
            class={cn(
              "min-w-0",
              horizontal() && "flex overflow-x-auto overscroll-x-contain pb-3",
            )}
            tabindex={horizontal() ? 0 : undefined}
          >
            <For each={props.items}>
              {(item, index) => (
                <li
                  class={cn(
                    "relative flex min-w-0 gap-3",
                    horizontal()
                      ? "min-w-48 flex-1 flex-col pr-6 last:pr-0"
                      : "pb-6 last:pb-0",
                  )}
                >
                  <Show when={index() < props.items.length - 1}>
                    <span
                      class={
                        horizontal()
                          ? "absolute top-4 right-0 left-8 border-t border-border"
                          : "absolute top-8 bottom-0 left-4 border-l border-border"
                      }
                      aria-hidden="true"
                    />
                  </Show>
                  <span
                    class={cn(
                      "relative flex size-8 shrink-0 items-center justify-center rounded-full border [&>svg]:size-4",
                      tones[item.tone ?? "default"],
                    )}
                    aria-hidden="true"
                  >
                    {item.icon ?? <CircleCheck />}
                  </span>
                  <div class="min-w-0 flex-1 pt-1">
                    <div
                      class={cn(
                        "flex min-w-0 flex-wrap items-start gap-x-4 gap-y-1",
                        horizontal() ? "flex-col" : "justify-between",
                      )}
                    >
                      <p class="text-sm font-medium wrap-anywhere">
                        {item.title}
                      </p>
                      <time
                        datetime={date(item)?.toISOString()}
                        class="text-xs text-muted-foreground"
                      >
                        {date(item)
                          ? (
                              props.dateFormatter ??
                              ((value) =>
                                new Intl.DateTimeFormat("pt-BR", {
                                  dateStyle: "short",
                                  timeStyle: "short",
                                }).format(value))
                            )(date(item)!)
                          : String(item.date)}
                      </time>
                    </div>
                    <Show when={item.description}>
                      <div class="mt-1 text-sm leading-relaxed text-muted-foreground wrap-anywhere">
                        {item.description}
                      </div>
                    </Show>
                    <Show when={item.action}>
                      <div class="mt-2 flex flex-wrap gap-2">{item.action}</div>
                    </Show>
                  </div>
                </li>
              )}
            </For>
          </ol>
        </Show>
      </Show>
    </div>
  );
}
