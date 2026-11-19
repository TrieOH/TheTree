import { Show, createMemo } from "solid-js";
import TrendingUp from "~icons/lucide/trending-up";
import TrendingDown from "~icons/lucide/trending-down";
import Minus from "~icons/lucide/minus";
import { cn } from "./lib/cn";

export interface MetricDeltaProps {
  value: number;
  previousValue: number;
  /** Whether an increase or a decrease represents an improvement. */
  improvement?: "increase" | "decrease";
  label?: string;
  class?: string;
}

/** Zero and invalid baselines never produce Infinity or a misleading percentage. */
export function MetricDelta(props: MetricDeltaProps) {
  const change = createMemo(() => {
    if (!Number.isFinite(props.value) || !Number.isFinite(props.previousValue))
      return null;
    if (props.previousValue === 0) return props.value === 0 ? 0 : null;
    return (
      ((props.value - props.previousValue) / Math.abs(props.previousValue)) *
      100
    );
  });
  const direction = () =>
    change() === null || change() === 0
      ? "flat"
      : change()! > 0
        ? "up"
        : "down";
  const positive = () =>
    direction() !== "flat" &&
    (direction() === "up") === (props.improvement !== "decrease");
  const formatted = () =>
    change() === null
      ? "Sem base de comparação"
      : new Intl.NumberFormat("pt-BR", {
          maximumFractionDigits: 1,
          signDisplay: "exceptZero",
        }).format(change()!) + "%";
  return (
    <div
      class={cn(
        "flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1 text-xs",
        props.class,
      )}
    >
      <span
        class={cn(
          "inline-flex items-center gap-1 font-medium tabular-nums",
          direction() === "flat"
            ? "text-muted-foreground"
            : positive()
              ? "text-emerald-700 dark:text-emerald-400"
              : "text-destructive",
        )}
      >
        <Show
          when={direction() === "up"}
          fallback={
            <Show
              when={direction() === "down"}
              fallback={<Minus class="size-3.5 shrink-0" aria-hidden="true" />}
            >
              <TrendingDown class="size-3.5 shrink-0" aria-hidden="true" />
            </Show>
          }
        >
          <TrendingUp class="size-3.5 shrink-0" aria-hidden="true" />
        </Show>
        {formatted()}
      </span>
      <span class="text-muted-foreground">
        {props.label ?? "vs. período anterior"}
      </span>
    </div>
  );
}
