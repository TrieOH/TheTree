import { shapeClass, type Shape } from "./lib/shape";
import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

export interface EmptyStateProps {
  shape?: Shape;
  /** Pass an icon element; this package deliberately carries no icon library. */
  icon?: JSX.Element;
  eyebrow?: JSX.Element;
  title: JSX.Element;
  description?: JSX.Element;
  action?: JSX.Element;
  class?: string;
  /** `inline` is a compact dashed box for an empty section inside a page. */
  variant?: "default" | "inline";
}

export function EmptyState(props: EmptyStateProps) {
  return (
    <Show
      when={props.variant !== "inline"}
      fallback={<InlineEmptyState {...props} />}
    >
      <div
        class={cn(
          "relative overflow-hidden rounded-2xl border border-border/60 bg-card px-6 py-8 text-center shadow-sm sm:px-8 sm:py-10",
          shapeClass(props.shape, "surface"),
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
    </Show>
  );
}

function InlineEmptyState(props: EmptyStateProps) {
  return (
    <div
      class={cn(
        "flex flex-col items-center gap-1 rounded-xl border border-dashed border-border px-4 py-6 text-center text-[13px] text-muted-foreground",
        shapeClass(props.shape, "surface"),
        props.class,
      )}
    >
      <Show when={props.icon}>
        <div class="mb-1 [&_svg]:size-4">{props.icon}</div>
      </Show>
      <p>{props.title}</p>
      <Show when={props.description}>
        <p>{props.description}</p>
      </Show>
      <Show when={props.action}>
        <div class="mt-2">{props.action}</div>
      </Show>
    </div>
  );
}
