import type { JSX } from "@solidjs/web";
import { For, Show } from "solid-js";
import { cn } from "./lib/cn";

export type MeterTone =
  | "default"
  | "accent"
  | "success"
  | "warning"
  | "destructive";

export interface MeterMarker {
  /** Position on the same scale as `value`. */
  value: number;
  /** Shown on hover. */
  label: string;
  tone?: "default" | "destructive";
}

export interface MeterProps {
  value: number;
  /** Defaults to 1, so `value` can be a share. */
  max?: number;
  /** Name of what is measured, on the left above the bar. */
  label?: JSX.Element;
  /** Reading on the right above the bar, like "3.2 of 10 GB". */
  valueLabel?: JSX.Element;
  /** Required when there is no text `label`. */
  ariaLabel?: string;
  tone?: MeterTone;
  /** Ticks across the bar, like a limit or an expected pace. */
  markers?: readonly MeterMarker[];
  class?: string;
}

const tones: Record<MeterTone, string> = {
  default: "bg-foreground/70",
  accent: "bg-accent",
  success: "bg-success",
  warning: "bg-warning",
  destructive: "bg-destructive dark:bg-destructive-foreground",
};

const markerTones = {
  default: "bg-foreground/60",
  destructive: "bg-destructive dark:bg-destructive-foreground",
};

/** How much of something is used, with optional ticks for limits. */
export function Meter(props: MeterProps) {
  const max = () => (props.max && props.max > 0 ? props.max : 1);
  const percent = (value: number) =>
    `${Math.min(100, Math.max(0, (value / max()) * 100))}%`;
  return (
    <div class={cn("space-y-1.5", props.class)}>
      <Show when={props.label || props.valueLabel}>
        <div class="flex items-baseline justify-between gap-3 text-[13px]">
          <span class="font-medium">{props.label}</span>
          <span class="tabular-nums text-muted-foreground">
            {props.valueLabel}
          </span>
        </div>
      </Show>
      <div
        role="meter"
        aria-label={
          props.ariaLabel ??
          (typeof props.label === "string" ? props.label : undefined)
        }
        aria-valuemin={0}
        aria-valuemax={max()}
        aria-valuenow={props.value}
        class="relative h-2 overflow-hidden rounded-full bg-muted"
      >
        <div
          class={cn("h-full rounded-full", tones[props.tone ?? "default"])}
          style={{ width: percent(props.value) }}
        />
        <For each={props.markers ?? []}>
          {(marker) => (
            <span
              aria-hidden="true"
              title={marker.label}
              class={cn(
                "absolute top-0 h-full w-0.5",
                markerTones[marker.tone ?? "default"],
              )}
              style={{ left: percent(marker.value) }}
            />
          )}
        </For>
      </div>
    </div>
  );
}
