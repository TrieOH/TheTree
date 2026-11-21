import { For, Show, createMemo } from "solid-js";
import { GenericChart } from "./generic-chart";
import { cn } from "./lib/cn";

export interface SparklineProps {
  data: readonly number[];
  label?: string;
  /** `bar` draws one bar per value; a zero is a dot on the baseline. */
  type?: "line" | "area" | "bar";
  /** Bar mode only: stretch the bars to fill the width, so strips of the same length line up. */
  fill?: boolean;
  color?: string;
  height?: number;
  showTooltip?: boolean;
  valueFormatter?: (value: number) => string;
  class?: string;
}
/** A compact chart without axes, using the same renderer as full-size charts. */
export function Sparkline(props: SparklineProps) {
  const data = createMemo(() =>
    props.data.filter(Number.isFinite).map((value, index) => ({
      date: new Date(Date.UTC(2000, 0, index + 1)),
      value,
      series: "sparkline",
    })),
  );
  return (
    <Show
      when={props.type !== "bar"}
      fallback={
        <SparkBars
          data={props.data}
          label={props.label}
          color={props.color}
          height={props.height}
          fill={props.fill}
          valueFormatter={props.valueFormatter}
          class={props.class}
        />
      }
    >
      <div class={cn("w-full min-w-0", props.class)}>
        <GenericChart
          data={data()}
          type={props.type ?? "line"}
          series={[
            {
              key: "sparkline",
              label: props.label ?? "Tendência",
              color: props.color ?? "var(--primary)",
            },
          ]}
          ariaLabel={props.label ?? "Histórico da métrica"}
          showAxes={false}
          showTooltip={props.showTooltip !== false}
          valueFormatter={props.valueFormatter}
          dateFormatter={(date) =>
            `Ponto ${Math.round((date.getTime() - Date.UTC(2000, 0, 1)) / 86_400_000) + 1}`
          }
          showPoints={data().length === 1}
          height={props.height ?? 56}
          emptyMessage="Sem histórico"
        />
      </div>
    </Show>
  );
}

function SparkBars(props: Omit<SparklineProps, "type" | "showTooltip">) {
  const values = createMemo(() => props.data.filter(Number.isFinite));
  const max = createMemo(() => Math.max(1, ...values()));
  const format = (value: number) =>
    props.valueFormatter?.(value) ?? String(value);
  return (
    <div
      role="img"
      aria-label={props.label ?? "Histórico da métrica"}
      class={cn("flex min-w-0 items-end gap-px", props.class)}
      style={{ height: `${props.height ?? 32}px` }}
    >
      <For each={values()}>
        {(value) => (
          <span
            title={format(value)}
            class={cn(
              props.fill ? "min-w-0 flex-1" : "w-1.5 shrink-0",
              "rounded-[2px]",
              value ? undefined : "h-px bg-border",
            )}
            style={
              value
                ? {
                    height: `${Math.max(8, (value / max()) * 100)}%`,
                    background: props.color ?? "var(--primary)",
                  }
                : undefined
            }
          />
        )}
      </For>
    </div>
  );
}
