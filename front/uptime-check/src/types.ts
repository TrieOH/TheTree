/** Status of a single 5-min check slot. Degraded = recovered on retry. */
export type PointStatus = "up" | "degraded" | "dead";

/** One raw 5-min data point. */
export interface Point {
  /** Epoch ms of the check start. */
  t: number;
  s: PointStatus;
  /** Round-trip ms of the successful attempt. */
  ms?: number;
  /** Error text of the last failed attempt (HTTP status or message). */
  err?: string;
}

/**
 * Persistent alert state per endpoint. Only "dead" and "up" are states —
 * degraded is a property of a data point, not a state (it stays silent).
 */
export interface EndpointState {
  status: "up" | "dead";
  /** Epoch ms the current state started. */
  since: number;
}

/** Uptime percentages per window, maintained by the daily compaction. */
export interface Rollup {
  "7d": number | null;
  "30d": number | null;
  "1y": number | null;
}

/** raw:<date> value — one KV key per UTC day, all endpoints inside. */
export type RawDay = Record<string, Point[]>;

/** hourly:<date> value — per endpoint, uptime % per UTC hour (null = no data). */
export type HourlyDay = Record<string, (number | null)[]>;

/** daily:<date> value — per endpoint, average of the day's hourly values. */
export type DailyDay = Record<string, number>;
