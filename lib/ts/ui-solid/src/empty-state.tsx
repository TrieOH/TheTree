import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

export interface EmptyStateProps {
  /** Pass an icon element; this package deliberately carries no icon library. */
  icon?: JSX.Element;
  eyebrow?: JSX.Element;
  title: JSX.Element;
  description?: JSX.Element;
  action?: JSX.Element;
  class?: string;
}

export function EmptyState(props: EmptyStateProps) {
  return (
    <div
      class={cn(
        "relative overflow-hidden rounded-2xl border border-border/60 bg-card px-6 py-8 text-center shadow-sm sm:px-8 sm:py-10",
        props.class,
      )}
    >
      <div class="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.08),transparent_42%)]" />

      <div class="relative z-10 mx-auto flex max-w-md flex-col items-center">
        <Show when={props.eyebrow}>
          <p class="mb-3 text-xs uppercase tracking-[0.24em] text-muted-foreground">
            {props.eyebrow}
          </p>
        </Show>

        <Show when={props.icon}>
          <div class="mb-5 flex h-14 w-14 items-center justify-center rounded-2xl border border-border/60 bg-muted/70 text-muted-foreground shadow-inner">
            {props.icon}
          </div>
        </Show>

        <h3 class="text-lg font-semibold tracking-tight text-foreground sm:text-xl">
          {props.title}
        </h3>

        <Show when={props.description}>
          <p class="mt-2 max-w-sm text-sm leading-relaxed text-muted-foreground sm:text-base">
            {props.description}
          </p>
        </Show>

        <Show when={props.action}>
          <div class="mt-5 flex justify-center">{props.action}</div>
        </Show>
      </div>
    </div>
  );
}
