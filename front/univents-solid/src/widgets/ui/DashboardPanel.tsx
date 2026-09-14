import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "@trieoh/ui-solid";

export interface DashboardPanelProps {
  title: string;
  description?: string;
  icon?: (props: { class?: string }) => JSX.Element;
  children: JSX.Element;
  class?: string;
}

export function DashboardPanel(props: DashboardPanelProps): JSX.Element {
  return (
    <section class={cn("min-w-0 max-w-full space-y-3", props.class)}>
      <div class="flex min-w-0 items-center gap-3 px-1">
        <Show when={props.icon}>
          {(iconFn) => (
            <div class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              {iconFn()({ class: "size-5" })}
            </div>
          )}
        </Show>
        <div class="min-w-0">
          <h2 class="text-base font-semibold tracking-tight text-foreground">
            {props.title}
          </h2>
          <Show when={props.description}>
            <p class="mt-0.5 truncate text-xs text-muted-foreground">
              {props.description}
            </p>
          </Show>
        </div>
      </div>
      {props.children}
    </section>
  );
}
