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
    const grouped = new Map<string, { date: string; revenue: number }>();
    for (const purchase of visible) {
      const date = new Date(purchase.timestamp);
      const key =
        range === "1h" || range === "24h"
          ? format(date, "yyyy-MM-dd HH:00")
          : format(date, "yyyy-MM-dd");
      const point = grouped.get(key) ?? {
        date:
          range === "1h" || range === "24h"
            ? format(date, "dd/MM HH:00")
            : format(date, "dd/MM"),
        revenue: 0,
      };
      if (purchase.status === "approved")
        point.revenue += purchase.totalCents / 100;
      grouped.set(key, point);
    }
    const ordered = [...grouped.entries()].sort(([a], [b]) =>
      a.localeCompare(b),
    );
    let accumulatedRevenue = 0;
    const points = ordered.map(([, point]) => {
      accumulatedRevenue += point.revenue;
      return { ...point, revenue: accumulatedRevenue };
    });
    if (points.length === 0) return points;
    return [
      {
        date: "",
        revenue: 0,
      },
      ...points,
    ];
  }, [range, visible]);
  const renderData = chartData.length ? chartData : [{ date: "", revenue: 0 }];
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
      axis: { label: "Período" },
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
    tooltip,
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
