import type { ServiceId } from "@/features/catalog/services";
import type { ColdStartDay } from "./model";

export interface ColdStartSummary {
  count: number;
  failed: number;
  /** Median start time, weighted by each day's starts. */
  typicalMs: number | null;
  /** The slowest day's p95. */
  worstP95Ms: number | null;
  /** Starts per day, oldest first, `days` long. */
  daily: number[];
}

/** One service's cold starts over the last `days` days (UTC). */
export function summarize(rows: readonly ColdStartDay[], service: ServiceId, days: number, now = new Date()): ColdStartSummary {
  const dayKeys = Array.from({ length: days }, (_, i) =>
    new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate() - (days - 1 - i))).toISOString().slice(0, 10),
  );
  const mine = rows.filter((row) => row.service === service && dayKeys.includes(row.day));
  const count = mine.reduce((n, row) => n + row.count, 0);
  return {
    count,
    failed: mine.reduce((n, row) => n + row.failed, 0),
    typicalMs: count ? mine.reduce((n, row) => n + row.p50Ms * row.count, 0) / count : null,
    worstP95Ms: mine.length ? Math.max(...mine.map((row) => row.p95Ms)) : null,
    daily: dayKeys.map((day) => mine.find((row) => row.day === day)?.count ?? 0),
  };
}
