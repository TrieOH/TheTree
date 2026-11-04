const compact = new Intl.NumberFormat("en", { notation: "compact", maximumFractionDigits: 1 });

export const formatCount = (value: number) => compact.format(value);

export function formatPercent(share: number | null, digits = 2): string {
  if (share === null) return "—";
  return `${(share * 100).toFixed(share === 1 || share === 0 ? 0 : digits)}%`;
}

export function formatMs(value: number | null): string {
  if (value === null) return "—";
  if (value < 1) return `${(value * 1000).toFixed(0)} µs`;
  if (value < 1000) return `${value.toFixed(value < 10 ? 1 : 0)} ms`;
  return `${(value / 1000).toFixed(1)} s`;
}

export function formatBytes(value: number | null): string {
  if (value === null) return "—";
  const units = ["B", "KB", "MB", "GB", "TB"];
  let n = value;
  let unit = 0;
  while (n >= 1024 && unit < units.length - 1) {
    n /= 1024;
    unit++;
  }
  return `${n.toFixed(n < 10 && unit > 0 ? 1 : 0)} ${units[unit]}`;
}

export function formatHours(seconds: number | null): string {
  if (seconds === null) return "—";
  const hours = seconds / 3600;
  return hours < 1 ? `${Math.round(seconds / 60)} min` : `${hours.toFixed(1)} h`;
}

export function formatUsage(value: number, unit: string): string {
  if (unit === "GiB-h") return `${value.toFixed(1)} GiB-h`;
  if (unit === "ms") return `${formatCount(value)} ms`;
  return formatCount(value);
}
