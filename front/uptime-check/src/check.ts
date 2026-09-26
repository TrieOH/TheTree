import type { Endpoint } from "./config";
import { CHECK, ENDPOINTS, NTFY } from "./config";
import { appendPoints, dayString, getState, putState } from "./kv";
import { notify } from "./ntfy";
import { probe } from "./probe";
import type { Point } from "./types";
import { formatDuration, sleep } from "./util";

/**
 * Run one 5-min check cycle for every endpoint.
 *
 * Flow per endpoint: initial probe; on failure, up to CHECK.retries retries
 * at CHECK.retryGapMs gaps (same invocation — scheduled handlers get 15 min
 * of wall time). Outcomes:
 *   - first attempt ok       -> up
 *   - recovered on a retry   -> degraded (silent, counts as up in averages)
 *   - all attempts failed    -> dead (urgent ntfy ping, once on transition)
 *
 * All endpoints run in parallel; results land in KV as ONE batched write
 * (KV free tier is ~1k writes/day — per-endpoint writes would blow it).
 */
export async function runChecks(env: Env): Promise<void> {
  console.log(
    `[runChecks] start urls=${ENDPOINTS.map((e) => e.url).join(", ")}`,
  );
  const entries = await Promise.all(
    ENDPOINTS.map(
      async (ep): Promise<readonly [string, Point]> => [
        ep.name,
        await checkEndpoint(env, ep),
      ],
    ),
  );

  const points: Record<string, Point> = {};
  for (const [name, point] of entries) {
    points[name] = point;
    console.log(
      `[runChecks] ${name}: ${point.s}${point.err ? ` (${point.err})` : ""}${point.ms ? ` ${point.ms}ms` : ""}`,
    );
  }

  await appendPoints(env.UPTIME, dayString(0), points);
  console.log(
    `[runChecks] done, ${Object.keys(points).length} points -> raw:${dayString(0)}`,
  );
}

/** Check one endpoint (with retries) and update its alert state. */
async function checkEndpoint(env: Env, ep: Endpoint): Promise<Point> {
  console.log(`[check:${ep.name}] probing ${ep.url}`);
  let last = await probe(ep.url);
  console.log(
    `[check:${ep.name}] initial: ${last.ok ? "ok" : `fail (${last.err})`} ${last.ms}ms`,
  );
  const firstFailed = !last.ok;
  if (firstFailed) {
    for (let i = 0; i < CHECK.retries; i++) {
      await sleep(CHECK.retryGapMs);
      last = await probe(ep.url);
      console.log(
        `[check:${ep.name}] retry ${i + 1}/${CHECK.retries}: ${last.ok ? "ok" : `fail (${last.err})`}`,
      );
      if (last.ok) break;
    }
  }

  const point: Point = last.ok
    ? firstFailed
      ? { t: last.t, s: "degraded", ms: last.ms }
      : { t: last.t, s: "up", ms: last.ms }
    : { t: last.t, s: "dead", err: last.err };

  console.log(
    `[check:${ep.name}] outcome=${point.s} firstFailed=${firstFailed}`,
  );
  await updateState(env, ep, point);
  return point;
}

/**
 * Alert on state transitions only: up->dead pings loudly, dead->up pings
 * a recovery with the downtime duration. Degraded never pings and never
 * touches the persistent state — it only marks the data point.
 * No transition -> no KV write (state writes are near-zero this way).
 */
async function updateState(
  env: Env,
  ep: Endpoint,
  point: Point,
): Promise<void> {
  const state = await getState(env.UPTIME, ep.name);
  console.log(
    `[state:${ep.name}] persistent=${state.status} since=${new Date(state.since).toISOString()}`,
  );

  if (point.s === "dead" && state.status !== "dead") {
    console.log(`[state:${ep.name}] TRANSITION -> dead, sending ntfy ping`);
    await notify(
      env,
      `${ep.name} is down`,
      `${point.err ?? "unreachable"} — retries exhausted (checked ${CHECK.retries + 1}x over ~2min).`,
      NTFY.deadPriority,
      [...NTFY.deadTags],
    );
    await putState(env.UPTIME, ep.name, { status: "dead", since: point.t });
  } else if (point.s !== "dead" && state.status === "dead") {
    console.log(`[state:${ep.name}] TRANSITION -> up, sending recovery ping`);
    await notify(
      env,
      `${ep.name} recovered`,
      `Back up after ${formatDuration(point.t - state.since)}.`,
      NTFY.recoveryPriority,
      [...NTFY.recoveryTags],
    );
    await putState(env.UPTIME, ep.name, { status: "up", since: point.t });
  }
}
