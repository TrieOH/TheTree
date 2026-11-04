import type { ServiceId } from "@/features/catalog/services";
import type { ServiceStatus } from "@/widgets/shell/StatusDot";

// Shared by the Worker (which pulls Cloudflare analytics and Neon into D1 and
// serves GET /api/metrics) and the pages. Nothing here ever calls a service.

/** What the Workers Paid plan includes per billing cycle. */
export const PLAN = {
  containers_memory: { label: "Container memory", unit: "GiB-h", included: 25 },
  workers_requests: { label: "Workers requests", unit: "requests", included: 10_000_000 },
  workers_cpu: { label: "Workers CPU time", unit: "ms", included: 30_000_000 },
  d1_rows_read: { label: "D1 rows read", unit: "rows", included: 25_000_000_000 },
  d1_rows_written: { label: "D1 rows written", unit: "rows", included: 50_000_000 },
  kv_reads: { label: "KV reads", unit: "reads", included: 10_000_000 },
  kv_writes: { label: "KV writes, deletes and lists", unit: "operations", included: 1_000_000 },
} as const;

export type PlanResource = keyof typeof PLAN;

/** Status over the last 24 hours. */
export const STATUS_WINDOW_HOURS = 24;
/** Degraded at this share of 5xx answers, once there are this many. */
export const DEGRADED_SHARE = 0.05;
export const DEGRADED_MIN_ERRORS = 3;

/**
 * Idle when nobody called the API in the window: a sleeping service is not a
 * broken one. "unknown" until the first pull has run.
 */
export function statusOf(requests: number | null, errors5xx: number | null): ServiceStatus {
  if (requests === null) return "unknown";
  if (requests === 0) return "idle";
  const errors = errors5xx ?? 0;
  return errors >= DEGRADED_MIN_ERRORS && errors / requests >= DEGRADED_SHARE ? "degraded" : "operational";
}

export interface ServiceMetrics {
  id: ServiceId;
  status: ServiceStatus;
  requests24h: number | null;
  errors24h: number | null;
  /** Share of requests not answered 5xx over the last 30 days; null without traffic. */
  successRate30d: number | null;
  /** Requests per hour, the last 24 hours, oldest first. */
  hourly: number[];
  lastRequestAt: string | null;
}

export interface Usage {
  resource: PlanResource;
  used: number;
  included: number;
  /** The guard's own cut-off, where one exists (containers stop at 90%). */
  stopsAt?: number;
}

export interface WorkerStats {
  script: string;
  service: ServiceId | null;
  kind: "front" | "api" | "portal" | "other";
  requests: number;
  errors: number;
  cpuP50Ms: number | null;
  cpuP99Ms: number | null;
  wallP50Ms: number | null;
  wallP99Ms: number | null;
}

export interface NeonStats {
  service: ServiceId;
  projectId: string;
  periodStart: string;
  periodEnd: string;
  storageBytes: number | null;
  activeSeconds: number | null;
  /** Share of the period the compute was suspended (scaled to zero). */
  suspendedShare: number | null;
  lastActiveAt: string | null;
  resumes: number;
  resumeP50Ms: number | null;
  resumeMaxMs: number | null;
}

export interface ColdStartDay {
  service: ServiceId;
  /** YYYY-MM-DD */
  day: string;
  count: number;
  failed: number;
  p50Ms: number;
  p95Ms: number;
}

export interface MetricsPayload {
  cycle: { start: string; end: string };
  services: ServiceMetrics[];
  usage: Usage[];
  workers: WorkerStats[];
  neon: NeonStats[];
  /** Container cold starts per service per day, the last 30 days, oldest first. */
  coldStarts: ColdStartDay[];
  /** When Cloudflare and Neon were last pulled; null before the first run. */
  syncedAt: string | null;
}

/** Billing cycles turn over on the 24th (UTC), as in @trieoh/container-budget. */
export function cycleBounds(now: Date, day = 24): { start: Date; end: Date } {
  const start = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), day));
  if (now.getUTCDate() < day) start.setUTCMonth(start.getUTCMonth() - 1);
  const end = new Date(start);
  end.setUTCMonth(end.getUTCMonth() + 1);
  return { start, end };
}
