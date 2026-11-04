import { COLD_START_DATASET, cycleStart, spentGiBSeconds } from "@trieoh/container-budget";
import { SERVICES } from "@/features/catalog/services";

// Every 30 minutes: plan usage, per-service traffic and per-Worker stats from
// Cloudflare's GraphQL Analytics API (read-only token), the container budget
// from the KV the guard writes, and each service's Neon project from Neon's
// control-plane API. None of these reach a container or Postgres.

const GRAPHQL = "https://api.cloudflare.com/client/v4/graphql";
const NEON = "https://console.neon.tech/api/v2";
const HOUR = 3_600_000;
/** Re-read the last few hours: analytics land a little late. */
const TRAFFIC_OVERLAP_HOURS = 3;
/** On the first run, how far back traffic is pulled. */
const TRAFFIC_BACKFILL_HOURS = 7 * 24;

async function graphql<T>(env: Env, query: string, variables: Record<string, unknown>): Promise<T> {
  const response = await fetch(GRAPHQL, {
    method: "POST",
    headers: { Authorization: `Bearer ${env.CLOUDFLARE_ANALYTICS_TOKEN}`, "Content-Type": "application/json" },
    body: JSON.stringify({ query, variables }),
  });
  const body = (await response.json()) as { data?: T; errors?: { message: string }[] | null };
  if (!response.ok || body.errors?.length || !body.data) {
    throw new Error(`Cloudflare analytics (${response.status}): ${body.errors?.map((e) => e.message).join("; ") ?? "no data"}`);
  }
  return body.data;
}

async function neon<T>(env: Env, path: string): Promise<T> {
  const response = await fetch(`${NEON}${path}`, {
    headers: { Authorization: `Bearer ${env.NEON_API_KEY}`, Accept: "application/json" },
  });
  if (!response.ok) throw new Error(`Neon ${path}: ${response.status}`);
  return response.json() as Promise<T>;
}

const hourOf = (date: Date) => new Date(Math.floor(date.getTime() / HOUR) * HOUR).toISOString().replace(".000Z", "Z");

// ── Traffic per API host, hourly ────────────────────────────────────────────

const TRAFFIC = `query($zone: string!, $since: Time!, $until: Time!, $hosts: [string!]) {
  viewer { zones(filter: { zoneTag: $zone }) {
    httpRequestsAdaptiveGroups(limit: 10000, filter: { datetime_geq: $since, datetime_lt: $until, clientRequestHTTPHost_in: $hosts }) {
      count
      dimensions { datetimeHour clientRequestHTTPHost edgeResponseStatus }
    } } } }`;

interface TrafficGroup {
  count: number;
  dimensions: { datetimeHour: string; clientRequestHTTPHost: string; edgeResponseStatus: number };
}

async function syncTraffic(env: Env, now: Date): Promise<D1PreparedStatement[]> {
  const last = await env.DB.prepare("SELECT MAX(bucket) AS bucket FROM metric_rollups WHERE metric = 'requests'")
    .first<{ bucket: string | null }>();
  const since = last?.bucket
    ? new Date(Date.parse(last.bucket) - TRAFFIC_OVERLAP_HOURS * HOUR)
    : new Date(now.getTime() - TRAFFIC_BACKFILL_HOURS * HOUR);

  const zones = new Map<string, string[]>();
  for (const service of SERVICES) {
    zones.set(service.analytics.zoneId, [...(zones.get(service.analytics.zoneId) ?? []), service.analytics.apiHost]);
  }
  const serviceOfHost = new Map(SERVICES.map((service) => [service.analytics.apiHost, service.id]));

  // Every hour in the window starts at zero, so a quiet hour is recorded as
  // quiet, not missing.
  const totals = new Map<string, { requests: number; errors: number }>();
  for (let hour = new Date(hourOf(since)); hour <= now; hour = new Date(hour.getTime() + HOUR)) {
    for (const service of SERVICES) totals.set(`${service.id}|${hourOf(hour)}`, { requests: 0, errors: 0 });
  }

  for (const [zone, hosts] of zones) {
    // The API caps one query's range at a day.
    for (let from = since; from < now; from = new Date(from.getTime() + 24 * HOUR)) {
      const until = new Date(Math.min(from.getTime() + 24 * HOUR, now.getTime()));
      const data = await graphql<{ viewer: { zones: { httpRequestsAdaptiveGroups: TrafficGroup[] }[] } }>(env, TRAFFIC, {
        zone,
        hosts,
        since: from.toISOString(),
        until: until.toISOString(),
      });
      for (const group of data.viewer.zones[0]?.httpRequestsAdaptiveGroups ?? []) {
        const service = serviceOfHost.get(group.dimensions.clientRequestHTTPHost);
        if (!service) continue;
        const key = `${service}|${group.dimensions.datetimeHour}`;
        const total = totals.get(key) ?? { requests: 0, errors: 0 };
        total.requests += group.count;
        if (group.dimensions.edgeResponseStatus >= 500) total.errors += group.count;
        totals.set(key, total);
      }
    }
  }

  const upsert = env.DB.prepare(
    `INSERT INTO metric_rollups (service, metric, bucket, value) VALUES (?, ?, ?, ?)
     ON CONFLICT (service, metric, bucket) DO UPDATE SET value = excluded.value`,
  );
  return [...totals].flatMap(([key, total]) => {
    const [service, bucket] = key.split("|");
    return [upsert.bind(service, "requests", bucket, total.requests), upsert.bind(service, "errors_5xx", bucket, total.errors)];
  });
}

// ── Plan usage and per-Worker stats, this billing cycle ─────────────────────

const ACCOUNT = `query($account: string!, $since: Time!, $until: Time!, $sinceDate: Date!) {
  viewer { accounts(filter: { accountTag: $account }) {
    workersInvocationsAdaptive(limit: 1000, filter: { datetime_geq: $since, datetime_lt: $until }) {
      sum { requests errors cpuTimeUs }
      quantiles { cpuTimeP50 cpuTimeP99 wallTimeP50 wallTimeP99 }
      dimensions { scriptName }
    }
    d1AnalyticsAdaptiveGroups(limit: 100, filter: { date_geq: $sinceDate }) {
      sum { rowsRead rowsWritten }
    }
    kvOperationsAdaptiveGroups(limit: 100, filter: { date_geq: $sinceDate }) {
      sum { requests }
      dimensions { actionType }
    }
  } } }`;

interface AccountData {
  viewer: {
    accounts: {
      workersInvocationsAdaptive: {
        sum: { requests: number; errors: number; cpuTimeUs: number };
        quantiles: { cpuTimeP50: number; cpuTimeP99: number; wallTimeP50: number; wallTimeP99: number };
        dimensions: { scriptName: string };
      }[];
      d1AnalyticsAdaptiveGroups: { sum: { rowsRead: number; rowsWritten: number } }[];
      kvOperationsAdaptiveGroups: { sum: { requests: number }; dimensions: { actionType: string } }[];
    }[];
  };
}

async function syncAccount(env: Env, now: Date): Promise<D1PreparedStatement[]> {
  const period = cycleStart(now);
  const data = await graphql<AccountData>(env, ACCOUNT, {
    account: env.CLOUDFLARE_ACCOUNT_ID,
    since: `${period}T00:00:00Z`,
    until: now.toISOString(),
    sinceDate: period,
  });
  const account = data.viewer.accounts[0];
  const workers = account?.workersInvocationsAdaptive ?? [];
  const sum = <T>(list: T[], pick: (item: T) => number) => list.reduce((total, item) => total + pick(item), 0);
  const kv = account?.kvOperationsAdaptiveGroups ?? [];
  const d1 = account?.d1AnalyticsAdaptiveGroups ?? [];
  const at = now.toISOString();

  const usage: [string, number, number][] = [
    ["containers_memory", (await spentGiBSeconds(env.BUDGET, now)) / 3600, 25],
    ["workers_requests", sum(workers, (w) => w.sum.requests), 10_000_000],
    ["workers_cpu", sum(workers, (w) => w.sum.cpuTimeUs) / 1000, 30_000_000],
    ["d1_rows_read", sum(d1, (g) => g.sum.rowsRead), 25_000_000_000],
    ["d1_rows_written", sum(d1, (g) => g.sum.rowsWritten), 50_000_000],
    ["kv_reads", sum(kv.filter((g) => g.dimensions.actionType === "read"), (g) => g.sum.requests), 10_000_000],
    ["kv_writes", sum(kv.filter((g) => g.dimensions.actionType !== "read"), (g) => g.sum.requests), 1_000_000],
  ];

  const usageUpsert = env.DB.prepare(
    `INSERT INTO usage_snapshots (resource, period, used, included, at) VALUES (?, ?, ?, ?, ?)
     ON CONFLICT (resource, period) DO UPDATE SET used = excluded.used, included = excluded.included, at = excluded.at`,
  );
  const workerUpsert = env.DB.prepare(
    `INSERT INTO worker_stats (script, period, requests, errors, cpu_p50_us, cpu_p99_us, wall_p50_us, wall_p99_us, at)
     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
     ON CONFLICT (script, period) DO UPDATE SET requests = excluded.requests, errors = excluded.errors,
       cpu_p50_us = excluded.cpu_p50_us, cpu_p99_us = excluded.cpu_p99_us,
       wall_p50_us = excluded.wall_p50_us, wall_p99_us = excluded.wall_p99_us, at = excluded.at`,
  );
  return [
    ...usage.map(([resource, used, included]) => usageUpsert.bind(resource, period, used, included, at)),
    ...workers
      .filter((w) => w.dimensions.scriptName && w.dimensions.scriptName !== "__unknown__")
      .map((w) =>
        workerUpsert.bind(
          w.dimensions.scriptName,
          period,
          w.sum.requests,
          w.sum.errors,
          w.quantiles.cpuTimeP50,
          w.quantiles.cpuTimeP99,
          w.quantiles.wallTimeP50,
          w.quantiles.wallTimeP99,
          at,
        ),
      ),
  ];
}

// ── Container cold starts, per service per day ─────────────────────────────

const ANALYTICS_ENGINE = (account: string) =>
  `https://api.cloudflare.com/client/v4/accounts/${account}/analytics_engine/sql`;

// Written by @trieoh/container-budget: blob1 the app, blob2 ok|failed,
// double1 milliseconds. _sample_interval undoes Analytics Engine's sampling.
const COLD_STARTS = `SELECT
    blob1 AS service,
    toStartOfDay(timestamp) AS day,
    SUM(_sample_interval) AS count,
    sumIf(_sample_interval, blob2 = 'failed') AS failed,
    quantileExactWeighted(0.5, double1, _sample_interval) AS p50,
    quantileExactWeighted(0.95, double1, _sample_interval) AS p95
  FROM ${COLD_START_DATASET}
  WHERE timestamp > NOW() - INTERVAL '30' DAY
  GROUP BY service, day
  FORMAT JSON`;

interface ColdStartRow {
  service: string;
  day: string;
  count: number | string;
  failed: number | string;
  p50: number | string;
  p95: number | string;
}

async function syncColdStarts(env: Env): Promise<D1PreparedStatement[]> {
  const response = await fetch(ANALYTICS_ENGINE(env.CLOUDFLARE_ACCOUNT_ID), {
    method: "POST",
    headers: { Authorization: `Bearer ${env.CLOUDFLARE_ANALYTICS_TOKEN}` },
    body: COLD_STARTS,
  });
  if (!response.ok) throw new Error(`Analytics Engine (${response.status}): ${(await response.text()).slice(0, 200)}`);
  const { data } = (await response.json()) as { data: ColdStartRow[] };
  const services = new Set(SERVICES.map((service) => service.id as string));
  const upsert = env.DB.prepare(
    `INSERT INTO cold_starts (resource, kind, day, count, p50_ms, p95_ms, failed) VALUES (?, 'container', ?, ?, ?, ?, ?)
     ON CONFLICT (resource, day) DO UPDATE SET count = excluded.count, p50_ms = excluded.p50_ms,
       p95_ms = excluded.p95_ms, failed = excluded.failed`,
  );
  return data
    .filter((row) => services.has(row.service))
    .map((row) =>
      upsert.bind(
        row.service,
        String(row.day).slice(0, 10),
        Math.round(Number(row.count)),
        Math.round(Number(row.p50)),
        Math.round(Number(row.p95)),
        Math.round(Number(row.failed)),
      ),
    );
}

// ── Neon, per service ───────────────────────────────────────────────────────

// GET /projects/{id} names these differently from the list endpoint
// (`active_time_seconds` there is `active_time` here), and carries the
// project's own consumption period.
interface NeonProject {
  project: {
    active_time_seconds?: number;
    compute_last_active_at?: string;
    consumption_period_start?: string;
    consumption_period_end?: string;
  };
}
interface NeonBranches {
  branches: { default?: boolean; primary?: boolean; logical_size?: number }[];
}
interface NeonOperations {
  operations: { action: string; status: string; created_at: string; total_duration_ms?: number }[];
  pagination?: { cursor?: string };
}

function median(values: number[]): number | null {
  if (!values.length) return null;
  const sorted = [...values].sort((a, b) => a - b);
  const mid = Math.floor(sorted.length / 2);
  return sorted.length % 2 ? sorted[mid] : Math.round((sorted[mid - 1] + sorted[mid]) / 2);
}

async function syncNeon(env: Env, now: Date): Promise<D1PreparedStatement[]> {
  const statements: D1PreparedStatement[] = [];
  for (const service of SERVICES) {
    const id = service.analytics.neonProject;
    const [{ project }, { branches }] = await Promise.all([
      neon<NeonProject>(env, `/projects/${id}`),
      neon<NeonBranches>(env, `/projects/${id}/branches`),
    ]);
    // Neon's own consumption period, which needn't start on the 1st.
    const end = project.consumption_period_end
      ? new Date(project.consumption_period_end)
      : new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth() + 1, 1));
    const start = project.consumption_period_start ? new Date(project.consumption_period_start) : new Date(end);
    if (!project.consumption_period_start) start.setUTCMonth(start.getUTCMonth() - 1);

    const resumes: number[] = [];
    let cursor: string | undefined;
    for (let page = 0; page < 20; page++) {
      const body = await neon<NeonOperations>(env, `/projects/${id}/operations?limit=100${cursor ? `&cursor=${encodeURIComponent(cursor)}` : ""}`);
      const inPeriod = body.operations.filter((op) => Date.parse(op.created_at) >= start.getTime());
      for (const op of inPeriod) {
        if (op.action === "start_compute" && op.status === "finished") resumes.push(op.total_duration_ms ?? 0);
      }
      cursor = body.pagination?.cursor;
      if (!cursor || inPeriod.length < body.operations.length || body.operations.length === 0) break;
    }

    const branch = branches.find((b) => b.default || b.primary) ?? branches[0];
    statements.push(
      env.DB.prepare(
        `INSERT INTO neon_projects (service, project_id, period_start, period_end, storage_bytes, active_seconds,
           last_active_at, resumes, resume_p50_ms, resume_max_ms, at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
         ON CONFLICT (service) DO UPDATE SET project_id = excluded.project_id, period_start = excluded.period_start,
           period_end = excluded.period_end, storage_bytes = excluded.storage_bytes, active_seconds = excluded.active_seconds,
           last_active_at = excluded.last_active_at, resumes = excluded.resumes, resume_p50_ms = excluded.resume_p50_ms,
           resume_max_ms = excluded.resume_max_ms, at = excluded.at`,
      ).bind(
        service.id,
        id,
        start.toISOString(),
        end.toISOString(),
        branch?.logical_size ?? null,
        project.active_time_seconds ?? null,
        project.compute_last_active_at ?? null,
        resumes.length,
        median(resumes),
        resumes.length ? Math.max(...resumes) : null,
        now.toISOString(),
      ),
    );
  }
  return statements;
}

export async function syncMetrics(env: Env, now = new Date()): Promise<{ statements: number }> {
  const [traffic, account, neonRows, coldStarts] = await Promise.all([
    syncTraffic(env, now),
    syncAccount(env, now),
    syncNeon(env, now),
    // The dataset only exists once an API Worker has written to it; until
    // then (or if it fails) the rest of the pull still lands.
    syncColdStarts(env).catch((error: unknown) => {
      console.warn("cold starts not pulled:", error instanceof Error ? error.message : error);
      return [];
    }),
  ]);
  const statements = [
    ...traffic,
    ...account,
    ...neonRows,
    ...coldStarts,
    env.DB.prepare(
      `INSERT INTO sync_state (source, synced_at) VALUES ('metrics', ?)
       ON CONFLICT (source) DO UPDATE SET synced_at = excluded.synced_at`,
    ).bind(now.toISOString()),
  ];
  for (let i = 0; i < statements.length; i += 50) await env.DB.batch(statements.slice(i, i + 50));
  return { statements: statements.length };
}

