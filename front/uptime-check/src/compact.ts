import { ENDPOINTS, RETENTION } from "./config";
import { dayString, getHourlyDay, getRawDay, keys } from "./kv";
import type { DailyDay, HourlyDay, Point, RawDay, Rollup } from "./types";
import { mean, pct1 } from "./util";

const DAY_S = 86_400;
const range = (start: number, end: number): number[] =>
  Array.from({ length: end - start }, (_, i) => start + i);

/**
 * Daily compaction (03:40 UTC cron). Rolls the retention ladder down one
 * tier and rebuilds the status-page rollups:
 *
 *   raw day 7        -> hourly:<date>  (avg uptime % per hour, degraded = up)
 *   hourly day 30    -> daily:<date>   (avg of the day's hourly values)
 *   daily day 365    -> deleted
 *
 * Each tier's day is compacted exactly once — one write + one delete per
 * tier per day, plus one rollup write per endpoint. Bounded by design.
 */
export async function runCompaction(env: Env): Promise<void> {
  await compactRawToHourly(env);
  await compactHourlyToDaily(env);
  await pruneDaily(env);
  await rebuildRollups(env);
}

async function compactRawToHourly(env: Env): Promise<void> {
  const day = dayString(RETENTION.rawDays);
  const raw: RawDay | null = await getRawDay(env.UPTIME, day);
  if (!raw) return;

  const hourly = hourlyFromRaw(raw);
  // 35d TTL is a safety net; compaction deletes the key after 30d anyway.
  await env.UPTIME.put(keys.hourly(day), JSON.stringify(hourly), {
    expirationTtl: (RETENTION.hourlyDays + 5) * DAY_S,
  });
  await env.UPTIME.delete(keys.raw(day));
  console.log(`compacted raw:${day} -> hourly`);
}

/** Per hour: share of non-dead checks, as % (null = no data that hour). */
function hourlyFromRaw(raw: RawDay): HourlyDay {
  const out: HourlyDay = {};
  for (const [name, points] of Object.entries(raw)) {
    const buckets: Point[][] = Array.from({ length: 24 }, () => []);
    for (const p of points) {
      buckets[new Date(p.t).getUTCHours()].push(p);
    }
    out[name] = buckets.map((b) =>
      b.length === 0
        ? null
        : pct1((b.filter((p) => p.s !== "dead").length / b.length) * 100),
    );
  }
  return out;
}

async function compactHourlyToDaily(env: Env): Promise<void> {
  const day = dayString(RETENTION.hourlyDays);
  const hourly: HourlyDay | null = await getHourlyDay(env.UPTIME, day);
  if (!hourly) return;

  const daily: DailyDay = {};
  for (const [name, hours] of Object.entries(hourly)) {
    const values = hours.filter((v): v is number => v !== null);
    if (values.length > 0) daily[name] = pct1(mean(values));
  }
  // 365d is KV's max TTL — the explicit prune handles anything older.
  await env.UPTIME.put(keys.daily(day), JSON.stringify(daily), {
    expirationTtl: RETENTION.dailyDays * DAY_S,
  });
  await env.UPTIME.delete(keys.hourly(day));
  console.log(`compacted hourly:${day} -> daily`);
}

async function pruneDaily(env: Env): Promise<void> {
  // Delete is a no-op when the key doesn't exist.
  await env.UPTIME.delete(keys.daily(dayString(RETENTION.dailyDays)));
}

/**
 * Rebuild per-endpoint uptime % for the status page. Every day in a window
 * contributes one number (raw day -> % of ok checks; hourly day -> mean of
 * its hours; daily day -> its stored value); the window is the mean of days.
 */
async function rebuildRollups(env: Env): Promise<void> {
  for (const ep of ENDPOINTS) {
    const [rawDays, hourlyDays, dailyDays] = await Promise.all([
      Promise.all(
        range(0, RETENTION.rawDays).map((d) =>
          getRawDay(env.UPTIME, dayString(d)),
        ),
      ),
      Promise.all(
        range(RETENTION.rawDays, RETENTION.hourlyDays).map((d) =>
          getHourlyDay(env.UPTIME, dayString(d)),
        ),
      ),
      Promise.all(
        range(RETENTION.hourlyDays, RETENTION.dailyDays).map((d) =>
          env.UPTIME.get<number>(keys.daily(dayString(d)), "json"),
        ),
      ),
    ]);

    const dayPct: (number | null)[] = [];
    for (const raw of rawDays) {
      const points = raw?.[ep.name] ?? [];
      dayPct.push(
        points.length === 0
          ? null
          : pct1(
              (points.filter((p) => p.s !== "dead").length / points.length) *
                100,
            ),
      );
    }
    for (const hourly of hourlyDays) {
      const hours = hourly?.[ep.name] ?? [];
      const values = hours.filter((v): v is number => v !== null);
      dayPct.push(values.length === 0 ? null : pct1(mean(values)));
    }
    for (const daily of dailyDays) {
      dayPct.push(daily ?? null);
    }

    const rollup: Rollup = {
      "7d": windowPct(dayPct, RETENTION.rawDays),
      "30d": windowPct(dayPct, RETENTION.hourlyDays),
      "1y": windowPct(dayPct, RETENTION.dailyDays),
    };
    await env.UPTIME.put(keys.rollup(ep.name), JSON.stringify(rollup));
    console.log(`rollup ${ep.name}: ${JSON.stringify(rollup)}`);
  }
}

function windowPct(dayPct: (number | null)[], days: number): number | null {
  const values = dayPct.slice(0, days).filter((v): v is number => v !== null);
  return values.length === 0 ? null : pct1(mean(values));
}
