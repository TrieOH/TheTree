import { clientOnly } from "@solidjs/web";
import { Loading } from "./loading";
import type { GenericChartProps } from "./chart/generic-chart";
export type { GenericChartProps } from "./chart/generic-chart";

// Keep the renderer and D3 out of applications that only use basic UI controls.
const ChartRenderer = clientOnly(() => import("./chart/generic-chart"), {
  lazy: true,
  export: "GenericChart",
});
export function GenericChart(props: GenericChartProps) {
  return (
    <ChartRenderer
      {...props}
      fallback={
        <div
          class="flex w-full min-w-0 items-center justify-center"
          style={{ height: `${Math.max(24, props.height ?? 320)}px` }}
        >
          <Loading size="sm" tone="muted" label="Carregando gráfico…" />
        </div>
      }
    />
  );
}
