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
  { value: "custom", label: "Personalizado" },
];
const durations: Record<string, number> = {
  "1h": 3_600_000,
  "24h": 86_400_000,
  "7d": 604_800_000,
  "30d": 2_592_000_000,
  all: Infinity,
};

const bucketTimestamp = (timestamp: string, range: string) => {
  const date = new Date(timestamp);
  if (range === "1h") date.setSeconds(0, 0);
  else if (range === "24h") date.setMinutes(0, 0, 0);
  else date.setHours(0, 0, 0, 0);
  return date.toISOString();
};

const customBucket = (start: string, end: string) => {
  const duration = new Date(end).getTime() - new Date(start).getTime();
  return duration <= 3_600_000 ? "1h" : duration <= 86_400_000 ? "24h" : "day";
};

export function DashboardLineChart({
  points,
  ranges = defaultRanges,
}: DashboardLineChartProps) {
  const [range, setRange] = useState("30d");
  const [customOpen, setCustomOpen] = useState(false);
  const [customStart, setCustomStart] = useState("");
  const [customEnd, setCustomEnd] = useState("");
  const customLabel =
    customStart && customEnd
      ? `${format(new Date(customStart), "dd/MM HH:mm")} – ${format(new Date(customEnd), "dd/MM HH:mm")}`
      : "Personalizado";
  const rangeOptions = ranges.map((option) =>
    option.value === "custom" ? { ...option, label: customLabel } : option,
  );
  const visible = points.filter((point) => {
    if (point.status !== "approved" && point.status !== "succeeded") {
      return false;
    }
    const time = new Date(point.timestamp).getTime();
    if (range === "custom") {
      const start = customStart ? new Date(customStart).getTime() : -Infinity;
      const end = customEnd ? new Date(customEnd).getTime() : Infinity;
      return !Number.isNaN(time) && time >= start && time <= end;
    }
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
      const date = bucketTimestamp(
        point.timestamp,
        range === "custom" ? customBucket(customStart, customEnd) : range,
      );
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
    const data = [...grouped.values()].map((point) => {
      totalCents += point.revenueCents;
      return {
        date: point.date,
        revenue: totalCents / 100,
        purchases: point.purchases,
      };
    });
    const currentDate = bucketTimestamp(
      range === "custom" && customEnd
        ? new Date(customEnd).toISOString()
        : new Date().toISOString(),
      range === "custom" ? customBucket(customStart, customEnd) : range,
    );
    if (data.at(-1)?.date !== currentDate) {
      data.push({ date: currentDate, revenue: totalCents / 100, purchases: 0 });
    }
    return data;
  }, [customEnd, customStart, range, visible]);
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
          options={rangeOptions}
          placeholder="Período"
          onChange={(value) => {
            setRange(value);
            if (value === "custom") setCustomOpen(true);
          }}
          className="w-64 max-w-full"
          triggerClassName="h-9 w-full text-sm"
        />
      </div>
      {customOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="w-full max-w-lg space-y-4 rounded-lg bg-card p-5 text-card-foreground shadow-xl ring-1 ring-foreground/10">
            <div>
              <h2 id="dashboard-chart-custom-range" className="font-semibold">
                Período personalizado
              </h2>
              <p className="text-xs text-muted-foreground">
                Escolha as datas que deseja visualizar.
              </p>
            </div>
            <div className="grid gap-3">
              <label className="grid gap-1 text-sm">
                Início
                <input
                  type="datetime-local"
                  value={customStart}
                  onChange={(event) => setCustomStart(event.target.value)}
                  className="h-10 w-full rounded-md border border-input bg-background px-2 text-sm"
                />
              </label>
              <label className="grid gap-1 text-sm">
                Fim
                <input
                  type="datetime-local"
                  value={customEnd}
                  onChange={(event) => setCustomEnd(event.target.value)}
                  className="h-10 w-full rounded-md border border-input bg-background px-2 text-sm"
                />
              </label>
            </div>
            <div className="flex justify-end gap-2">
              <button
                type="button"
                className="h-9 rounded-md px-3 text-sm hover:bg-muted"
                onClick={() => {
                  setCustomOpen(false);
                  setRange("30d");
                }}
              >
                Cancelar
              </button>
              <button
                type="button"
                className="h-9 rounded-md bg-primary px-3 text-sm text-primary-foreground hover:bg-primary/90"
                disabled={!customStart || !customEnd || customStart > customEnd}
                onClick={() => setCustomOpen(false)}
              >
                Aplicar
              </button>
            </div>
          </div>
        </div>
      )}
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
