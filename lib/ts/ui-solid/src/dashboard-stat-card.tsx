import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { MetricDelta, type MetricDeltaProps } from "./metric-delta";
import { Sparkline, type SparklineProps } from "./sparkline";
import { Skeleton } from "./skeleton";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";

export interface DashboardStatCardProps {
  label: string;
  value: JSX.Element;
  hint?: string;
  icon?: (props: { class?: string }) => JSX.Element;
  action?: JSX.Element;
  comparison?: MetricDeltaProps;
  sparkline?: SparklineProps;
  loading?: boolean;
  shape?: Shape;
  class?: string;
}
export function DashboardStatCard(props: DashboardStatCardProps) {
  return (
    <article
      aria-label={props.label}
      aria-busy={props.loading ? "true" : "false"}
      class={cn(
        "min-w-0 rounded-xl border border-border bg-card p-4 shadow-xs",
        shapeClass(props.shape, "surface"),
        props.class,
      )}
    >
      <div class="flex flex-wrap items-center justify-between gap-2 text-muted-foreground">
        <p class="min-w-0 text-xs wrap-anywhere">{props.label}</p>
        <div class="flex shrink-0 items-center gap-2">
          {props.action}
          <Show when={props.icon}>
            <span>{props.icon?.({ class: "size-4" })}</span>
          </Show>
        </div>
      </div>
      <div class="mt-2 min-h-7 min-w-0 text-xl font-semibold tracking-tight text-foreground tabular-nums wrap-anywhere">
        <Show
          when={!props.loading}
          fallback={<Skeleton class="h-7 w-32 max-w-full" />}
        >
          {props.value}
        </Show>
      </div>
      <Show when={props.hint}>
        <p class="mt-2 text-xs leading-relaxed text-muted-foreground wrap-anywhere">
          {props.hint}
        </p>
      </Show>
      <Show when={!props.loading && props.comparison}>
        <MetricDelta
          value={props.comparison!.value}
          previousValue={props.comparison!.previousValue}
          improvement={props.comparison!.improvement}
          label={props.comparison!.label}
          class={cn("mt-3", props.comparison!.class)}
        />
      </Show>
      <Show when={!props.loading && props.sparkline}>
        <Sparkline
          data={props.sparkline!.data}
          label={props.sparkline!.label}
          type={props.sparkline!.type}
          color={props.sparkline!.color}
          height={props.sparkline!.height}
          showTooltip={props.sparkline!.showTooltip}
          valueFormatter={props.sparkline!.valueFormatter}
          class={cn("mt-3", props.sparkline!.class)}
        />
      </Show>
    </article>
  );
}
