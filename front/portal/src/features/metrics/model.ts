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
 * Also degraded when the container's p95 over the window is above this, once
 * there are this many requests. Measured inside the container, so cold
 * starts don't count.
 */
export const DEGRADED_P95_MS = 1500;
export const DEGRADED_MIN_REQUESTS_FOR_LATENCY = 20;

/**
 * Idle when nobody called the API in the window: a sleeping service is not a
 * broken one. "unknown" until the first pull has run.
 */
export function statusOf(requests: number | null, errors5xx: number | null, p95Ms: number | null = null): ServiceStatus {
  if (requests === null) return "unknown";
  if (requests === 0) return "idle";
  const errors = errors5xx ?? 0;
  if (errors >= DEGRADED_MIN_ERRORS && errors / requests >= DEGRADED_SHARE) return "degraded";
  if (p95Ms !== null && p95Ms > DEGRADED_P95_MS && requests >= DEGRADED_MIN_REQUESTS_FOR_LATENCY) return "degraded";
  return "operational";
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
  /** The container's p95 over the 24 hours before the last pull (Grafana); null without data. */
  p95Ms24h: number | null;
  /** p95 per hour, the last 24 hours, oldest first; null for an hour without requests. */
  p95Hourly: (number | null)[];
}

/** One endpoint over the last 7 days, from the containers' own metrics. */
export interface EndpointStats {
  service: ServiceId;
  method: string;
  /** The route template, e.g. /auth/login; never a concrete path. */
  route: string;
  requests: number;
  /** The worst day's p95 in the window. */
  p95Ms: number | null;
}

/** Problem answers by type slug over the last 7 days: slugs and counts only. */
export interface ProblemStats {
  service: ServiceId;
  problem: string;
  status: number;
  count: number;
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
  /** Slowest first. Empty until Grafana has been pulled. */
  endpoints: EndpointStats[];
  /** Most frequent first. */
  problems: ProblemStats[];
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
