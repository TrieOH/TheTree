import { addDays, isSameDay, startOfDay } from "./date-utils";
import type { ChartDatum } from "./types";

export type ContinuityMode = "continuous" | "none";

interface FillContinuousSeriesOptions {
  data: ChartDatum[];
  seriesKeys: string[];
  windowStart: Date | null;
  windowEnd: Date | null;
  mode: ContinuityMode;
}

export function fillContinuousSeries({
  data,
  seriesKeys,
  windowStart,
  windowEnd,
  mode,
}: FillContinuousSeriesOptions): ChartDatum[] {
  if (mode === "none" || data.length === 0) return data;

  const bySeries = new Map<string, ChartDatum[]>();
  for (const key of seriesKeys) bySeries.set(key, []);
  for (const datum of data) {
    bySeries.get(datum.series)?.push(datum);
  }

  const result: ChartDatum[] = [];

  for (const key of seriesKeys) {
    const points = (bySeries.get(key) ?? [])
      .slice()
      .sort((a, b) => a.date.getTime() - b.date.getTime());
    if (points.length === 0) continue;

    const from = windowStart ?? addDays(points[0].date, -1);
    const to = windowEnd ?? points[points.length - 1].date;
    if (from.getTime() > to.getTime()) continue;

    let pointIndex = 0;
    let lastKnownValue: number | null = null;
    const cursor = startOfDay(from);
    const end = startOfDay(to);

    while (
      pointIndex < points.length &&
      startOfDay(points[pointIndex].date).getTime() < cursor.getTime()
    ) {
      lastKnownValue = points[pointIndex].value;
      pointIndex++;
    }

    while (cursor.getTime() <= end.getTime()) {
      let matchedToday = false;
      while (
        pointIndex < points.length &&
        isSameDay(points[pointIndex].date, cursor)
      ) {
        lastKnownValue = points[pointIndex].value;
        result.push(points[pointIndex]);
        pointIndex++;
        matchedToday = true;
      }
      if (!matchedToday) {
        result.push({
          date: new Date(cursor),
          value: lastKnownValue ?? 0,
          series: key,
          synthetic: true,
        });
      }
      cursor.setDate(cursor.getDate() + 1);
    }
  }

  return result;
}
