import { previewUsage } from "./preview-usage";
import type { Values } from "../catalog";
import exampleSource from "../examples/DashboardExamples.tsx?raw";

const examples: Record<string, string> = {
  DashboardStatCard: "StatCardExample",
  DashboardPanel: "PanelExample",
  DashboardBarList: "BarListExample",
  MetricDelta: "MetricDeltaExample",
  Sparkline: "SparklineExample",
  GenericChart: "ChartExample",
  ChartCard: "ChartExample",
};

/** Generate copyable code from the actual preview, so its content and behavior stay aligned. */
export function dashboardUsage(
  name: string,
  values: Values,
): string | undefined {
  const example = examples[name];
  if (!example) return undefined;
  return previewUsage(
    exampleSource,
    example,
    values,
    exampleSource.slice(
      exampleSource.indexOf("const history ="),
      exampleSource.indexOf("export function StatCardExample"),
    ),
    name === "GenericChart",
  );
}
