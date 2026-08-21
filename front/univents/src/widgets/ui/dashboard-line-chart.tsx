import { defineChart, lineY } from "@tanstack/charts";
import { Chart } from "@tanstack/charts/react";
import { scaleLinear } from "@tanstack/charts/scales/linear";
import { scalePoint } from "@tanstack/charts/scales/point";
import { tooltip } from "@tanstack/charts/tooltip";
import { format } from "date-fns";
import { useMemo, useState } from "react";
import { ToolbarCombobox } from "@/features/certifications/editor/ui/toolbar-combobox";

export interface DashboardPurchasePoint {
  timestamp: string;
  status: string;
  totalCents: number;
}

type Range = "1h" | "24h" | "7d" | "30d" | "all";
const ranges: Array<{ value: Range; label: string }> = [
  { value: "1h", label: "1h" },
  { value: "24h", label: "24h" },
  { value: "7d", label: "7 dias" },
  { value: "30d", label: "30 dias" },
  { value: "all", label: "Tudo" },
];

interface DashboardLineChartProps {
  purchases: DashboardPurchasePoint[];
}

export function DashboardLineChart({ purchases }: DashboardLineChartProps) {
  const [range, setRange] = useState<Range>("30d");
  const now = Date.now();
  const milliseconds = {
    "1h": 3_600_000,
    "24h": 86_400_000,
    "7d": 604_800_000,
    "30d": 2_592_000_000,
    all: Infinity,
  }[range];
  const visible = purchases.filter((purchase) => {
    const time = new Date(purchase.timestamp).getTime();
    return !Number.isNaN(time) && now - time <= milliseconds;
  });
  const chartData = useMemo(() => {
    const ordered = [...visible].sort(
      (a, b) =>
        new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
    );
    let accumulatedRevenue = 0;
    const points = ordered.map((purchase) => {
      const date = new Date(purchase.timestamp);
      if (purchase.status === "approved") {
        accumulatedRevenue += purchase.totalCents / 100;
      }
      return {
        date: date.toISOString(),
        revenue: accumulatedRevenue,
      };
    });
    return points;
  }, [range, visible]);
  const renderData = chartData.length ? chartData : [{ date: "", revenue: 0 }];
  const tickLabels = new Map<string, string>();
  const renderedDays = new Set<string>();
  for (const point of renderData) {
    const date = new Date(point.date);
    if (Number.isNaN(date.getTime())) continue;
    const day = format(date, "yyyy-MM-dd");
    tickLabels.set(
      point.date,
      renderedDays.has(day) ? "" : format(date, "dd/MM"),
    );
    renderedDays.add(day);
  }
  const definition = defineChart({
    marks: [
      lineY(renderData, {
        id: "revenue",
        x: "date",
        y: "revenue",
        points: chartData.length > 0,
        stroke: "#10b981",
        strokeWidth: 2.5,
        strokeOpacity: chartData.length > 0 ? 1 : 0,
      }),
    ],
    x: {
      scale: () => scalePoint<string>().padding(0),
      axis: {
        label: "Período",
        ticks: {
          format: (value) => tickLabels.get(value) ?? "",
        },
      },
    },
    y: {
      scale: () =>
        scaleLinear().domain([
          0,
          Math.max(...renderData.map((point) => point.revenue), 1),
        ]),
      nice: true,
      grid: true,
      axis: { label: "Lucro (R$)" },
    },
    tooltip: {
      use: tooltip,
      format: (point) => {
        const date = new Date(String(point.xValue));
        const period = Number.isNaN(date.getTime())
          ? "Período desconhecido"
          : format(date, "dd/MM HH:mm");
        return `Período: ${period}\nLucro acumulado: R$ ${Number(point.yValue).toFixed(2)}`;
      },
    },
  });

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap items-center justify-end gap-2">
        <ToolbarCombobox
          value={range}
          options={ranges}
          placeholder="Período"
          searchPlaceholder="Filtrar período…"
          onChange={(value) => setRange(value as Range)}
          className="w-40"
          triggerClassName="h-9 w-full text-sm"
        />
      </div>
      <div className="min-w-0 max-w-full overflow-hidden outline-none focus-within:outline-none focus-within:ring-0 [&_*:focus]:outline-none [&_*:focus]:ring-0">
        <Chart
          definition={definition}
          height={240}
          initialWidth={640}
          ariaLabel="Lucro acumulado ao longo do tempo"
          ariaDescription="Crescimento do valor acumulado das compras aprovadas no período selecionado."
        />
      </div>
    </div>
  );
}
