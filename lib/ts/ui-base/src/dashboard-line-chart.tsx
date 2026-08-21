import { areaY, d3Curve, defineChart, lineY } from "@tanstack/charts";
import { Chart } from "@tanstack/charts/react";
import { scaleLinear } from "@tanstack/charts/scales/linear";
import { scalePoint } from "@tanstack/charts/scales/point";
import { tooltip } from "@tanstack/charts/tooltip";
import { format } from "date-fns";
import { curveMonotoneX } from "d3-shape";
import { useMemo, useState } from "react";
import {
  ToolbarCombobox,
  type ToolbarComboboxOption,
} from "./toolbar-combobox";

export interface DashboardLineChartPoint {
  timestamp: string;
  status: string;
  totalCents: number;
}

export interface DashboardLineChartProps {
  points: DashboardLineChartPoint[];
  ranges?: readonly ToolbarComboboxOption[];
}

const defaultRanges = [
  { value: "1h", label: "1 hora" },
  { value: "24h", label: "24 horas" },
  { value: "7d", label: "7 dias" },
  { value: "30d", label: "30 dias" },
  { value: "all", label: "Tudo" },
];
const durations: Record<string, number> = {
  "1h": 3_600_000,
  "24h": 86_400_000,
  "7d": 604_800_000,
  "30d": 2_592_000_000,
  all: Infinity,
};

export function DashboardLineChart({
  points,
  ranges = defaultRanges,
}: DashboardLineChartProps) {
  const [range, setRange] = useState("30d");
  const visible = points.filter((point) => {
    const time = new Date(point.timestamp).getTime();
    return (
      !Number.isNaN(time) &&
      Date.now() - time <= (durations[range] ?? durations["30d"])
    );
  });
  const data = useMemo(() => {
    const grouped = new Map<
      string,
      { date: string; revenueCents: number; purchases: number }
    >();
    for (const point of [...visible].sort((a, b) =>
      a.timestamp.localeCompare(b.timestamp),
    )) {
      const dateValue = new Date(point.timestamp);
      dateValue.setSeconds(0, 0);
      const date = dateValue.toISOString();
      const current = grouped.get(date) ?? {
        date,
        revenueCents: 0,
        purchases: 0,
      };
      current.purchases += 1;
      if (point.status === "approved" || point.status === "succeeded") {
        current.revenueCents += point.totalCents;
      }
      grouped.set(date, current);
    }

    let totalCents = 0;
    return [...grouped.values()].map((point) => {
      totalCents += point.revenueCents;
      return {
        date: point.date,
        revenue: totalCents / 100,
        purchases: point.purchases,
      };
    });
  }, [visible]);
  const renderData = data.length
    ? data
    : [{ date: "", revenue: 0, purchases: 0 }];
  const labels = new Map<string, string>();
  const days = new Set<string>();
  for (const point of renderData) {
    const date = new Date(point.date);
    if (Number.isNaN(date.getTime())) continue;
    const day = format(date, "yyyy-MM-dd");
    labels.set(point.date, days.has(day) ? "" : format(date, "dd/MM"));
    days.add(day);
  }
  const definition = defineChart({
    marks: [
      areaY(renderData, {
        id: "revenue-area",
        key: "date",
        x: "date",
        y: "revenue",
        fill: "#10b981",
        fillOpacity: data.length > 0 ? 0.12 : 0,
        curve: d3Curve(curveMonotoneX),
      }),
      lineY(renderData, {
        id: "revenue",
        key: "date",
        x: "date",
        y: "revenue",
        points: data.length > 0,
        stroke: "#10b981",
        strokeWidth: 2.5,
        curve: d3Curve(curveMonotoneX),
        strokeOpacity: data.length > 0 ? 1 : 0,
      }),
    ],
    x: {
      scale: () => scalePoint<string>().padding(0),
      axis: {
        label: "Período",
        ticks: { format: (value) => labels.get(value) ?? "" },
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
      axis: {
        label: "Lucro (R$)",
        ticks: {
          format: (value) =>
            `R$ ${Number(value).toLocaleString("pt-BR", {
              maximumFractionDigits: 0,
            })}`,
        },
      },
    },
    tooltip: {
      use: tooltip,
      format: (point) => {
        const date = new Date(String(point.xValue));
        const purchases =
          typeof point.datum?.purchases === "number"
            ? point.datum.purchases
            : 0;
        return `Período: ${Number.isNaN(date.getTime()) ? "desconhecido" : format(date, "dd/MM HH:mm")}\nCompras: ${purchases}\nLucro acumulado: R$ ${Number(point.yValue).toFixed(2)}`;
      },
    },
  });
  return (
    <div className="space-y-3">
      <div className="flex justify-end">
        <ToolbarCombobox
          value={range}
          options={ranges}
          placeholder="Período"
          onChange={setRange}
          className="w-40"
          triggerClassName="h-9 w-full text-sm"
        />
      </div>
      <div className="min-w-0 max-w-full overflow-hidden">
        <Chart
          definition={definition}
          height={240}
          initialWidth={640}
          ariaLabel="Lucro acumulado ao longo do tempo"
        />
      </div>
    </div>
  );
}
