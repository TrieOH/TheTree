import type { EndpointState, HourlyDay, Point, RawDay } from "./types";

const pad = (n: number) => String(n).padStart(2, "0");

/** YYYY-MM-DD in UTC. */
export function dateKey(d: Date): string {
  return `${d.getUTCFullYear()}-${pad(d.getUTCMonth() + 1)}-${pad(d.getUTCDate())}`;
}

/** The UTC day string N days ago (0 = today). */
export function dayString(daysAgo: number): string {
  const d = new Date();
  d.setUTCDate(d.getUTCDate() - daysAgo);
  return dateKey(d);
}

/**
 * Key layout — designed so the 5-min cron does ONE write per run
 * (KV free tier allows ~1k writes/day; per-endpoint writes would cost ~2.3k):
 *
 *   raw:<date>            Record<endpoint, Point[]>   5-min points, 7d, TTL 8d
 *   hourly:<date>         Record<endpoint, (num|null)[]>  24 uptime %/hour, 30d
 *   daily:<date>          Record<endpoint, number>    avg of hourly, 1y
 *   state:<endpoint>      EndpointState               written on transitions only
 *   rollup:<endpoint>     Rollup                      7d/30d/1y %, rebuilt daily
 */
export const keys = {
  raw: (day: string) => `raw:${day}`,
  hourly: (day: string) => `hourly:${day}`,
  daily: (day: string) => `daily:${day}`,
  state: (name: string) => `state:${name}`,
  rollup: (name: string) => `rollup:${name}`,
};

/** Read one raw day (all endpoints), or null when absent/empty. */
export async function getRawDay(
  kv: KVNamespace,
  day: string,
): Promise<RawDay | null> {
  const value = await kv.get<RawDay>(keys.raw(day), "json");
  return value ?? null;
}

/**
 * Append points to today's raw day key. One read + one write per run,
 * regardless of how many endpoints were checked.
 */
export async function appendPoints(
  kv: KVNamespace,
  day: string,
  points: Record<string, Point>,
): Promise<void> {
  const existing = (await getRawDay(kv, day)) ?? {};
  console.log(
    `[kv] appendPoints raw:${day} had ${Object.values(existing).reduce((a, l) => a + l.length, 0)} points`,
  );
  for (const [name, point] of Object.entries(points)) {
    const list = existing[name] ?? [];
    list.push(point);
    existing[name] = list;
  }
  // 8d TTL is a safety net — compaction deletes the key after 7d anyway.
  await kv.put(keys.raw(day), JSON.stringify(existing), {
    expirationTtl: 8 * 24 * 3600,
  });
  console.log(
    `[kv] appendPoints raw:${day} wrote, now ${Object.values(existing).reduce((a, l) => a + l.length, 0)} points`,
  );
}

export async function getState(
  kv: KVNamespace,
  name: string,
): Promise<EndpointState> {
  const state = await kv.get<EndpointState>(keys.state(name), "json");
  return state ?? { status: "up", since: Date.now() };
}

/** Only called on transitions (rare) — keeps write volume near zero. */
export async function putState(
  kv: KVNamespace,
  name: string,
  state: EndpointState,
): Promise<void> {
  await kv.put(keys.state(name), JSON.stringify(state));
}

export async function getHourlyDay(
  kv: KVNamespace,
  day: string,
): Promise<HourlyDay | null> {
  const value = await kv.get<HourlyDay>(keys.hourly(day), "json");
  return value ?? null;
}
