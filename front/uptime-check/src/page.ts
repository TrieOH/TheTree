import { ENDPOINTS } from "./config";
import { dayString, getRawDay, getState, keys } from "./kv";
import type { EndpointState, Point, RawDay, Rollup } from "./types";

const COLORS = {
  up: "#22c55e",
  degraded: "#f59e0b",
  dead: "#ef4444",
  missing: "#2b3442",
} as const;

const esc = (s: string): string =>
  s
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");

/**
 * The public status page, served from KV — no box involved. Reads are
 * bounded: 7 raw days + state + rollup per page load (~15 KV reads).
 */
export async function renderPage(env: Env): Promise<Response> {
  const [rawDays, states, rollups] = await Promise.all([
    Promise.all(range(0, 7).map((d) => getRawDay(env.UPTIME, dayString(d)))),
    Promise.all(ENDPOINTS.map((ep) => getState(env.UPTIME, ep.name))),
    Promise.all(
      ENDPOINTS.map((ep) =>
        env.UPTIME.get<Rollup>(keys.rollup(ep.name), "json"),
      ),
    ),
  ]);

  const rows = ENDPOINTS.map((ep, i) =>
    rowHtml(ep.name, rawDays, states[i], rollups[i]),
  ).join("\n");

  const anyDead = states.some((s) => s.status === "dead");
  const anyDegraded = ENDPOINTS.some(
    (ep) => lastPoint(rawDays, ep.name)?.s === "degraded",
  );
  const banner = anyDead
    ? { text: "Partial outage", color: COLORS.dead }
    : anyDegraded
      ? { text: "Degraded performance", color: COLORS.degraded }
      : { text: "All systems operational", color: COLORS.up };

  const html = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>TheTree status</title>
<style>
  :root { color-scheme: dark; }
  body { margin: 0; background: #0b0f14; color: #e5e7eb;
         font: 15px/1.5 system-ui, -apple-system, sans-serif; }
  main { max-width: 720px; margin: 0 auto; padding: 48px 20px; }
  h1 { font-size: 20px; font-weight: 600; margin: 0; }
  .banner { display: flex; align-items: center; gap: 10px;
            background: #111827; border: 1px solid #1f2937; border-radius: 10px;
            padding: 16px 20px; margin: 24px 0; font-weight: 600; }
  .banner .dot { width: 12px; height: 12px; border-radius: 50%; }
  .row { background: #111827; border: 1px solid #1f2937; border-radius: 10px;
         padding: 16px 20px; margin-bottom: 12px; }
  .row .head { display: flex; justify-content: space-between; align-items: baseline; }
  .row .name { font-weight: 600; }
  .badge { font-size: 13px; font-weight: 600; }
  .strip { display: flex; gap: 2px; height: 28px; margin: 12px 0 6px; }
  .strip div { flex: 1 1 0; border-radius: 2px; min-width: 1px; }
  .uptimes { display: flex; gap: 16px; font-size: 13px; color: #9ca3af; }
  .uptimes b { color: #e5e7eb; font-weight: 600; }
  footer { margin-top: 28px; font-size: 13px; color: #6b7280; }
</style>
</head>
<body>
<main>
  <h1>TheTree status</h1>
  <div class="banner"><span class="dot" style="background:${banner.color}"></span>${banner.text}</div>
  ${rows}
  <footer>Checked every 5 minutes from Cloudflare's edge — outside the box. Strip shows the last 7 days at 5-minute resolution.</footer>
</main>
</body>
</html>`;

  return new Response(html, {
    headers: {
      "content-type": "text/html; charset=utf-8",
      "cache-control": "public, max-age=60",
    },
  });
}

function rowHtml(
  name: string,
  rawDays: (RawDay | null)[],
  state: EndpointState,
  rollup: Rollup | null,
): string {
  const last = lastPoint(rawDays, name);
  const badge =
    state.status === "dead"
      ? { text: "Down", color: COLORS.dead }
      : last?.s === "degraded"
        ? { text: "Degraded", color: COLORS.degraded }
        : { text: "Operational", color: COLORS.up };

  const blocks: string[] = [];
  for (const raw of rawDays) {
    const points = raw?.[name];
    if (!points || points.length === 0) {
      blocks.push(
        `<div style="background:${COLORS.missing}" title="no data"></div>`,
      );
      continue;
    }
    for (const p of points) {
      blocks.push(
        `<div style="background:${COLORS[p.s]}" title="${esc(title(p))}"></div>`,
      );
    }
  }

  const up = (v: number | null | undefined) =>
    v === null || v === undefined ? "—" : `${v}%`;

  return `<div class="row">
  <div class="head">
    <span class="name">${esc(name)}</span>
    <span class="badge" style="color:${badge.color}">${badge.text}</span>
  </div>
  <div class="strip">${blocks.join("")}</div>
  <div class="uptimes">
    <span>7d <b>${up(rollup?.["7d"])}</b></span>
    <span>30d <b>${up(rollup?.["30d"])}</b></span>
    <span>1y <b>${up(rollup?.["1y"])}</b></span>
  </div>
</div>`;
}

function title(p: Point): string {
  const time = `${new Date(p.t).toISOString().replace("T", " ").slice(0, 16)} UTC`;
  const base = `${time} — ${p.s}${p.ms ? ` ${p.ms}ms` : ""}`;
  return p.err ? `${base} (${p.err})` : base;
}

function lastPoint(
  rawDays: (RawDay | null)[],
  name: string,
): Point | undefined {
  const all = rawDays.flatMap((r) => r?.[name] ?? []);
  return all.reduce<Point | undefined>(
    (latest, p) => (!latest || p.t > latest.t ? p : latest),
    undefined,
  );
}

function range(start: number, end: number): number[] {
  return Array.from({ length: end - start }, (_, i) => start + i);
}
