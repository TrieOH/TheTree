import { createEffect, createMemo, createSignal, untrack } from "solid-js";
import { computeRangeWindow } from "./date-range-modal";
import type { ChartDatum, ChartType, CustomRange, RangeKey } from "./types";

interface UseChartFiltersOptions {
  data: () => ChartDatum[];
  seriesKeys: () => string[];
  initialType?: ChartType;
  initialRange?: RangeKey;
}

export function useChartFilters(options: UseChartFiltersOptions) {
  const [type, setType] = createSignal<ChartType>(options.initialType ?? "line");
  const [range, setRange] = createSignal<RangeKey>(options.initialRange ?? "all");
  const [customRange, setCustomRange] = createSignal<CustomRange>({
    from: null,
    to: null,
  });
  const [visible, setVisible] = createSignal<Set<string>>(
    new Set(untrack(() => options.seriesKeys())),
  );
  const [query, setQuery] = createSignal("");

  createEffect(
    () => options.seriesKeys(),
    (keys) => {
      setVisible((current) =>
        current.size === 0 && keys.length > 0 ? new Set(keys) : current,
      );
    },
  );

  const toggleSeries = (key: string) => {
    setVisible((prev) => {
      if (prev.has(key)) {
        if (prev.size === 1) return prev;
        const next = new Set(prev);
        next.delete(key);
        return next;
      }
      return new Set(prev).add(key);
    });
  };

  const isolateSeries = (key: string) => {
    setVisible(new Set([key]));
  };

  const showAllSeries = () => {
    setVisible(new Set(options.seriesKeys()));
  };

  const window = createMemo(() => computeRangeWindow(range(), customRange()));
  const windowStart = () => window().windowStart;
  const windowEnd = () => window().windowEnd;

  const seriesFilteredData = createMemo(() => {
    const list: ChartDatum[] = [];
    const q = query().trim().toLowerCase();
    const vis = visible();
    for (const d of options.data()) {
      if (vis.has(d.series) && (!q || d.series.toLowerCase().includes(q))) {
        list.push(d);
      }
    }
    return list;
  });

  const filteredData = createMemo(() => {
    const list: ChartDatum[] = [];
    const start = windowStart();
    const end = windowEnd();
    for (const d of seriesFilteredData()) {
      const t = d.date.getTime();
      if (start !== null && t < start.getTime()) continue;
      if (end !== null && t > end.getTime()) continue;
      list.push(d);
    }
    return list;
  });

  return {
    type,
    setType,
    range,
    setRange,
    customRange,
    setCustomRange,
    visible,
    toggleSeries,
    isolateSeries,
    showAllSeries,
    query,
    setQuery,
    windowStart,
    windowEnd,
    seriesFilteredData,
    filteredData,
  };
}
