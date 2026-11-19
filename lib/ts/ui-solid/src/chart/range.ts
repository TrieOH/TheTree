import { startOfDay } from "./date-utils";
import type { CustomRange, RangeKey, RangeOption } from "./types";
export const PRESETS: RangeOption[] = [
  { key: "7d", label: "Últimos 7 dias", days: 7 },
  { key: "30d", label: "Últimos 30 dias", days: 30 },
  { key: "90d", label: "Últimos 90 dias", days: 90 },
  { key: "all", label: "Tudo", days: null },
];
const endOfDay = (date: Date) => {
  const result = new Date(date);
  result.setHours(23, 59, 59, 999);
  return result;
};
export function computeRangeWindow(
  range: RangeKey,
  customRange: CustomRange,
  now = new Date(),
): { windowStart: Date | null; windowEnd: Date | null } {
  if (range === "all") return { windowStart: null, windowEnd: null };
  if (range === "custom") {
    if (!customRange.from) return { windowStart: null, windowEnd: null };
    const a = customRange.from,
      b = customRange.to ?? a;
    return {
      windowStart: startOfDay(a.getTime() <= b.getTime() ? a : b),
      windowEnd: endOfDay(a.getTime() <= b.getTime() ? b : a),
    };
  }
  const days = PRESETS.find((preset) => preset.key === range)?.days;
  if (!days) return { windowStart: null, windowEnd: null };
  const start = startOfDay(now);
  start.setDate(start.getDate() - days + 1);
  return { windowStart: start, windowEnd: endOfDay(now) };
}
