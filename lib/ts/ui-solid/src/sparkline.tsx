import { createMemo } from "solid-js";
import { GenericChart } from "./generic-chart";
import { cn } from "./lib/cn";

export interface SparklineProps {
  data: readonly number[];
  label?: string;
  type?: "line" | "area";
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
  );
}
