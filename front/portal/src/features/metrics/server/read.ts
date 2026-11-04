import { LIMIT_GIB_SECONDS, cycleStart } from "@trieoh/container-budget";
import { SERVICES, type ServiceId } from "@/features/catalog/services";
import {
  PLAN,
  STATUS_WINDOW_HOURS,
  statusOf,
  type MetricsPayload,
  type ColdStartDay,
  type EndpointStats,
  type NeonStats,
  type PlanResource,
  type ProblemStats,
  type ServiceMetrics,
  type Usage,
  type WorkerStats,
} from "../model";

// GET /api/metrics: everything the sync stored, shaped for the pages.

const HOUR = 3_600_000;

interface Rollup {
  service: ServiceId;
  metric: "requests" | "errors_5xx" | "p95_ms" | "p95_ms_24h";
  bucket: string;
  value: number;
}

export async function readMetrics(db: D1Database, now = new Date()): Promise<MetricsPayload> {
  const period = cycleStart(now);
  const end = new Date(`${period}T00:00:00Z`);
  end.setUTCMonth(end.getUTCMonth() + 1);
  const since30d = new Date(now.getTime() - 30 * 24 * HOUR).toISOString();
  const since7d = new Date(now.getTime() - 7 * 24 * HOUR).toISOString().slice(0, 10);

  const [rollups, usage, workers, neon, sync, cold, endpointRows, problemRows] = await db.batch<unknown>([
    db.prepare(
      "SELECT service, metric, bucket, value FROM metric_rollups WHERE metric IN ('requests', 'errors_5xx', 'p95_ms', 'p95_ms_24h') AND bucket >= ? ORDER BY bucket",
    ).bind(since30d),
    db.prepare("SELECT resource, used, included FROM usage_snapshots WHERE period = ?").bind(period),
    db.prepare(
      "SELECT script, requests, errors, cpu_p50_us, cpu_p99_us, wall_p50_us, wall_p99_us FROM worker_stats WHERE period = ? ORDER BY requests DESC",
    ).bind(period),
    db.prepare("SELECT * FROM neon_projects"),
    db.prepare("SELECT synced_at FROM sync_state WHERE source = 'metrics'"),
    db.prepare(
      "SELECT resource, day, count, failed, p50_ms, p95_ms FROM cold_starts WHERE kind = 'container' AND day >= ? ORDER BY day",
    ).bind(since30d.slice(0, 10)),
    db.prepare(
      `SELECT service, method, route, SUM(requests) AS requests, MAX(p95_ms) AS p95_ms FROM endpoint_rollups
       WHERE day >= ? GROUP BY service, method, route ORDER BY p95_ms DESC NULLS LAST, requests DESC LIMIT 50`,
    ).bind(since7d),
    db.prepare(
      `SELECT service, problem, MAX(status) AS status, SUM(count) AS count FROM problem_rollups
       WHERE day >= ? GROUP BY service, problem ORDER BY count DESC LIMIT 50`,
    ).bind(since7d),
  ]);

  const rows = rollups.results as Rollup[];
  const windowStart = new Date(now.getTime() - STATUS_WINDOW_HOURS * HOUR);
  const synced = (sync.results as { synced_at: string }[])[0]?.synced_at ?? null;

  const services = SERVICES.map((service): ServiceMetrics => {
    const mine = rows.filter((row) => row.service === service.id);
    const total = (metric: Rollup["metric"], from: Date) =>
      mine.filter((row) => row.metric === metric && Date.parse(row.bucket) >= from.getTime()).reduce((n, row) => n + row.value, 0);
    const hasData = mine.some((row) => row.metric === "requests");
    const requests24h = hasData ? total("requests", windowStart) : null;
    const errors24h = hasData ? total("errors_5xx", windowStart) : null;
    const requests30d = total("requests", new Date(since30d));
    const buckets = Array.from({ length: STATUS_WINDOW_HOURS }, (_, i) =>
      new Date(Math.floor(now.getTime() / HOUR) * HOUR - (STATUS_WINDOW_HOURS - 1 - i) * HOUR)
        .toISOString()
        .replace(".000Z", "Z"),
    );
    const at = (metric: Rollup["metric"], bucket: string) =>
      mine.find((row) => row.metric === metric && row.bucket === bucket)?.value;
    const hourly = buckets.map((bucket) => at("requests", bucket) ?? 0);
    const last = [...mine].reverse().find((row) => row.metric === "requests" && row.value > 0);
    // Only a recent pull's 24h p95 speaks for the status window.
    const p95 = [...mine].reverse().find((row) => row.metric === "p95_ms_24h");
    const p95Ms24h = p95 && Date.parse(p95.bucket) >= windowStart.getTime() ? p95.value : null;
    return {
      id: service.id,
      status: statusOf(requests24h, errors24h, p95Ms24h),
      p95Ms24h,
      p95Hourly: buckets.map((bucket) => at("p95_ms", bucket) ?? null),
      requests24h,
      errors24h,
      successRate30d: requests30d ? 1 - total("errors_5xx", new Date(since30d)) / requests30d : null,
      hourly,
      lastRequestAt: last?.bucket ?? null,
    };
  });

  const usageRows = usage.results as { resource: PlanResource; used: number; included: number }[];
  const usageOut = (Object.keys(PLAN) as PlanResource[]).flatMap((resource): Usage[] => {
    const row = usageRows.find((u) => u.resource === resource);
    if (!row) return [];
    return [
      {
        resource,
        used: row.used,
        included: row.included,
        ...(resource === "containers_memory" ? { stopsAt: LIMIT_GIB_SECONDS / 3600 } : {}),
      },
    ];
  });

  const byWorker = new Map<string, { service: ServiceId; kind: "front" | "api" }>();
  for (const service of SERVICES) {
    byWorker.set(service.analytics.frontWorker, { service: service.id, kind: "front" });
    byWorker.set(service.analytics.apiWorker, { service: service.id, kind: "api" });
  }
  const ms = (us: number | null) => (us === null ? null : us / 1000);
  const workersOut = (
    workers.results as {
      script: string;
      requests: number;
      errors: number;
      cpu_p50_us: number | null;
      cpu_p99_us: number | null;
      wall_p50_us: number | null;
      wall_p99_us: number | null;
    }[]
  ).map(
    (row): WorkerStats => ({
      script: row.script,
      service: byWorker.get(row.script)?.service ?? null,
      kind: byWorker.get(row.script)?.kind ?? (row.script === "portal" ? "portal" : "other"),
      requests: row.requests,
      errors: row.errors,
      cpuP50Ms: ms(row.cpu_p50_us),
      cpuP99Ms: ms(row.cpu_p99_us),
      wallP50Ms: ms(row.wall_p50_us),
      wallP99Ms: ms(row.wall_p99_us),
    }),
  );

  const neonOut = (
    neon.results as {
      service: ServiceId;
      project_id: string;
      period_start: string;
      period_end: string;
      storage_bytes: number | null;
      active_seconds: number | null;
      last_active_at: string | null;
      resumes: number;
      resume_p50_ms: number | null;
      resume_max_ms: number | null;
      at: string;
    }[]
  ).map((row): NeonStats => {
    const elapsed = (Math.min(Date.parse(row.at), Date.parse(row.period_end)) - Date.parse(row.period_start)) / 1000;
    return {
      service: row.service,
      projectId: row.project_id,
      periodStart: row.period_start,
      periodEnd: row.period_end,
      storageBytes: row.storage_bytes,
      activeSeconds: row.active_seconds,
      suspendedShare: row.active_seconds !== null && elapsed > 0 ? Math.max(0, 1 - row.active_seconds / elapsed) : null,
      lastActiveAt: row.last_active_at,
      resumes: row.resumes,
      resumeP50Ms: row.resume_p50_ms,
      resumeMaxMs: row.resume_max_ms,
    };
  });

  const coldStarts = (
    cold.results as { resource: ServiceId; day: string; count: number; failed: number; p50_ms: number; p95_ms: number }[]
  ).map(
    (row): ColdStartDay => ({
      service: row.resource,
      day: row.day,
      count: row.count,
      failed: row.failed,
      p50Ms: row.p50_ms,
      p95Ms: row.p95_ms,
    }),
  );

  const endpoints = (
    endpointRows.results as { service: ServiceId; method: string; route: string; requests: number; p95_ms: number | null }[]
  ).map((row): EndpointStats => ({ service: row.service, method: row.method, route: row.route, requests: row.requests, p95Ms: row.p95_ms }));
  const problems = (problemRows.results as { service: ServiceId; problem: string; status: number; count: number }[]).map(
    (row): ProblemStats => ({ service: row.service, problem: row.problem, status: row.status, count: row.count }),
  );

  return {
    cycle: { start: `${period}T00:00:00Z`, end: end.toISOString() },
    coldStarts,
    endpoints,
    problems,
    services,
    usage: usageOut,
    workers: workersOut,
    neon: neonOut,
    syncedAt: synced,
  };
}
