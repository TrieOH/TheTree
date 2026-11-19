import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { Button } from "./button";
import { Loading } from "./loading";
import { cn } from "./lib/cn";

export interface DashboardPanelProps {
  title: string;
  description?: string;
  icon?: (props: { class?: string }) => JSX.Element;
  action?: JSX.Element;
  children: JSX.Element;
  loading?: boolean;
  error?: string;
  onRetry?: () => void;
  class?: string;
}
export function DashboardPanel(props: DashboardPanelProps) {
  return (
    <section
      aria-label={props.title}
      aria-busy={props.loading ? "true" : "false"}
      class={cn("min-w-0 max-w-full space-y-3", props.class)}
    >
      <header class="flex min-w-0 flex-wrap items-start justify-between gap-3 px-1">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <Show when={props.icon}>
            <div class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              {props.icon?.({ class: "size-5" })}
            </div>
          </Show>
          <div class="min-w-0">
            <h2 class="text-base font-semibold tracking-tight wrap-anywhere">
              {props.title}
            </h2>
            <Show when={props.description}>
              <p class="mt-0.5 text-xs text-muted-foreground wrap-anywhere">
                {props.description}
              </p>
            </Show>
          </div>
        </div>
        <Show when={props.action}>
          <div class="ml-auto shrink-0">{props.action}</div>
        </Show>
      </header>
      <Show
        when={!props.error}
        fallback={
          <div class="rounded-md border border-border p-6 text-center">
            <p role="alert" class="text-sm text-muted-foreground">
              {props.error}
            </p>
            <Show when={props.onRetry}>
              <Button
                size="sm"
                variant="outline"
                class="mt-3"
                onClick={() => props.onRetry?.()}
              >
                Tentar novamente
              </Button>
            </Show>
          </div>
        }
      >
        <Show
          when={!props.loading}
          fallback={<Loading size="sm" tone="muted" class="min-h-32" />}
        >
          {props.children}
        </Show>
      </Show>
    </section>
  );
}
