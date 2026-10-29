import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { cn } from "./lib/cn";

export type SeparatorOrientation = "horizontal" | "vertical";
export type SeparatorVariant = "solid" | "dashed" | "dotted";
export type SeparatorLabelPosition = "start" | "center" | "end";

export interface SeparatorProps {
  /** Direction of the line. A vertical separator fills the height of its flex/grid parent. */
  orientation?: SeparatorOrientation;
  /** Line style. */
  variant?: SeparatorVariant;
  /** Text or element shown inside the line, e.g. "ou". */
  label?: JSX.Element;
  /** Where the label sits along the line. */
  labelPosition?: SeparatorLabelPosition;
  /**
   * `true` when the line is only visual: assistive tech skips it. Set `false`
   * when it divides distinct sections of content and should be announced.
   */
  decorative?: boolean;
  class?: string;
}

const lineStyle: Record<SeparatorVariant, string> = {
  solid: "border-solid",
  dashed: "border-dashed",
  dotted: "border-dotted",
};

/**
 * Hairline between groups of content, optionally carrying a label. Built from
 * borders (not a filled bar) so `dashed` and `dotted` work in both directions.
 */
export function Separator(props: SeparatorProps): JSX.Element {
  const orientation = () => props.orientation ?? "horizontal";
  const variant = () => props.variant ?? "solid";
  const labelPosition = () => props.labelPosition ?? "center";
  const decorative = () => props.decorative ?? true;
  const vertical = () => orientation() === "vertical";

  // A decorative separator is `role="none"`; a semantic one says which way it runs.
  const a11y = () =>
    decorative()
      ? { role: "none" as const }
      : {
          role: "separator" as const,
          "aria-orientation": orientation(),
        };

  const line = (extra?: string) => (
    <span
      aria-hidden="true"
      class={cn(
        "border-border",
        lineStyle[variant()],
        vertical() ? "border-l" : "border-t",
        extra,
      )}
    />
  );

  return (
    <Show
      when={props.label}
      fallback={
        <div
          {...a11y()}
          class={cn(
            "shrink-0 border-border",
            lineStyle[variant()],
            vertical() ? "self-stretch border-l" : "w-full border-t",
            props.class,
          )}
        />
      }
    >
      <div
        {...a11y()}
        class={cn(
          "flex shrink-0 items-center gap-3 text-xs text-muted-foreground",
          vertical() ? "flex-col self-stretch" : "w-full",
          props.class,
        )}
      >
        <Show when={labelPosition() !== "start"}>
          {line("flex-1")}
        </Show>
        <span class="shrink-0">{props.label}</span>
        <Show when={labelPosition() !== "end"}>
          {line("flex-1")}
        </Show>
      </div>
    </Show>
  );
}
