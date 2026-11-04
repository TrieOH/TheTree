import { SERVICES, type ServiceId } from "@/features/catalog/services";

// Latency, slowest endpoints and error types from Grafana Cloud's Prometheus
// query API, with a read-only access policy token. What's queried is what
// the API containers push over OTLP (lib/go/telemetry); the containers
// report as `<service>-container`. Nothing here reaches a container.

const HOUR = 3_600_000;
const DAY = 24 * HOUR;
/** Re-read the last few hours: OTLP batches land a little late. */
const OVERLAP_HOURS = 3;
/** On the first run, how far back hourly p95 is pulled. */
const BACKFILL_HOURS = 7 * 24;
const CONTAINERS = `service_name=~".+-container"`;

type Vector = { metric: Record<string, string>; value: [number, string] }[];
type Matrix = { metric: Record<string, string>; values: [number, string][] }[];

async function prom<T>(env: Env, endpoint: "query" | "query_range", params: Record<string, string>): Promise<T> {
  const response = await fetch(`${env.GRAFANA_PROM_URL}/api/v1/${endpoint}?${new URLSearchParams(params)}`, {
    headers: { Authorization: `Basic ${btoa(`${env.GRAFANA_PROM_USER}:${env.GRAFANA_READ_TOKEN}`)}` },
  });
  const body = (await response.json().catch(() => null)) as { status?: string; data?: { result: T }; error?: string } | null;
  if (!response.ok || body?.status !== "success" || !body.data) {
    throw new Error(`Grafana ${endpoint} (${response.status}): ${body?.error ?? "no data"}`);
  }
  return body.data.result;
}

const KNOWN = new Set<string>(SERVICES.map((service) => service.id));

/** `identityx-container` → `identityx`; anything else is not ours. */
function serviceOf(metric: Record<string, string>): ServiceId | null {
  const id = metric.service_name?.replace(/-container$/, "");
  return id && id !== metric.service_name && KNOWN.has(id) ? (id as ServiceId) : null;
}

const finite = (raw: string) => {
  const n = Number(raw);
  return Number.isFinite(n) ? n : null;
};

const iso = (ms: number) => new Date(ms).toISOString().replace(".000Z", "Z");

const p95 = (bucket: string, by: string, window: string) =>
  `histogram_quantile(0.95, sum by (le, ${by}) (increase(${bucket}{${CONTAINERS}}[${window}])))`;

async function syncLatency(env: Env, now: Date): Promise<D1PreparedStatement[]> {
  const last = await env.DB.prepare("SELECT MAX(bucket) AS bucket FROM metric_rollups WHERE metric = 'p95_ms'")
    .first<{ bucket: string | null }>();
  const hourNow = Math.floor(now.getTime() / HOUR) * HOUR;
  const since = last?.bucket ? Date.parse(last.bucket) - OVERLAP_HOURS * HOUR : hourNow - BACKFILL_HOURS * HOUR;

  // Evaluated at each full hour, the value covers the hour before it.
  const [hourly, day] = await Promise.all([
    prom<Matrix>(env, "query_range", {
      query: p95("http_server_request_duration_seconds_bucket", "service_name", "1h"),
      start: String((since + HOUR) / 1000),
      end: String(hourNow / 1000),
      step: "3600",
    }),
    prom<Vector>(env, "query", {
      query: p95("http_server_request_duration_seconds_bucket", "service_name", "24h"),
      time: String(now.getTime() / 1000),
    }),
  ]);

  const upsert = env.DB.prepare(
    `INSERT INTO metric_rollups (service, metric, bucket, value) VALUES (?, ?, ?, ?)
     ON CONFLICT (service, metric, bucket) DO UPDATE SET value = excluded.value`,
  );
  const statements: D1PreparedStatement[] = [];
  for (const series of hourly) {
    const service = serviceOf(series.metric);
    if (!service) continue;
    for (const [t, raw] of series.values) {
      const seconds = finite(raw);
      if (seconds !== null) statements.push(upsert.bind(service, "p95_ms", iso(t * 1000 - HOUR), seconds * 1000));
    }
  }
  for (const sample of day) {
    const service = serviceOf(sample.metric);
    const seconds = finite(sample.value[1]);
    if (service && seconds !== null) statements.push(upsert.bind(service, "p95_ms_24h", iso(hourNow), seconds * 1000));
  }
  return statements;
}

/** Yesterday (closed) and today so far, as UTC days. */
function days(now: Date): { day: string; at: number; window: number }[] {
  const today = Math.floor(now.getTime() / DAY) * DAY;
  return [
    { day: iso(today - DAY).slice(0, 10), at: today, window: DAY / 1000 },
    { day: iso(today).slice(0, 10), at: now.getTime(), window: Math.floor((now.getTime() - today) / 1000) },
  ].filter((d) => d.window >= 60);
}

async function syncEndpointsAndProblems(env: Env, now: Date): Promise<D1PreparedStatement[]> {
  const endpointUpsert = env.DB.prepare(
    `INSERT INTO endpoint_rollups (service, method, route, day, requests, p95_ms) VALUES (?, ?, ?, ?, ?, ?)
     ON CONFLICT (service, method, route, day) DO UPDATE SET requests = excluded.requests, p95_ms = excluded.p95_ms`,
  );
  const problemUpsert = env.DB.prepare(
    `INSERT INTO problem_rollups (service, problem, day, status, count) VALUES (?, ?, ?, ?, ?)
     ON CONFLICT (service, problem, day) DO UPDATE SET status = excluded.status, count = excluded.count`,
  );
  const statements: D1PreparedStatement[] = [];
  for (const { day, at, window } of days(now)) {
    const time = String(at / 1000);
    const [requests, latency, problems] = await Promise.all([
      prom<Vector>(env, "query", {
        query: `sum by (service_name, method, route) (increase(http_requests_total{${CONTAINERS}}[${window}s]))`,
        time,
      }),
      prom<Vector>(env, "query", { query: p95("http_request_duration_seconds_bucket", "service_name, route", `${window}s`), time }),
      prom<Vector>(env, "query", {
        query: `sum by (service_name, problem, status) (increase(http_problems_total{${CONTAINERS}}[${window}s]))`,
        time,
      }),
    ]);

    const p95ByRoute = new Map<string, number>();
    for (const sample of latency) {
      const seconds = finite(sample.value[1]);
      if (seconds !== null) p95ByRoute.set(`${sample.metric.service_name}|${sample.metric.route}`, seconds * 1000);
    }
    for (const sample of requests) {
      const service = serviceOf(sample.metric);
      const count = Math.round(finite(sample.value[1]) ?? 0);
      // Unmatched paths aren't endpoints; they'd crowd the table.
      if (!service || count < 1 || !sample.metric.route || sample.metric.route === "not_found") continue;
      const routeP95 = p95ByRoute.get(`${sample.metric.service_name}|${sample.metric.route}`) ?? null;
      statements.push(endpointUpsert.bind(service, sample.metric.method ?? "", sample.metric.route, day, count, routeP95));
    }
    for (const sample of problems) {
      const service = serviceOf(sample.metric);
      const count = Math.round(finite(sample.value[1]) ?? 0);
      if (!service || count < 1 || !sample.metric.problem) continue;
      statements.push(problemUpsert.bind(service, sample.metric.problem, day, Number(sample.metric.status) || 0, count));
    }
  }
  return statements;
}

/** Everything from Grafana; empty while the token or URLs aren't configured. */
export async function syncGrafana(env: Env, now: Date): Promise<D1PreparedStatement[]> {
  if (!env.GRAFANA_READ_TOKEN || !env.GRAFANA_PROM_URL || !env.GRAFANA_PROM_USER) return [];
  const [latency, daily] = await Promise.all([syncLatency(env, now), syncEndpointsAndProblems(env, now)]);
  return [...latency, ...daily];
}
