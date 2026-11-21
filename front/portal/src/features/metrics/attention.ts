import { findService, type ServiceId } from "@/features/catalog/services";
import { summarize } from "./cold-starts";
import { formatCount, formatMs, formatPercent } from "./format";
import { DEGRADED_P95_MS, PLAN, type MetricsPayload, type Usage } from "./model";

/** Which Metrics tab explains an item. */
export type MetricsTab = "services" | "latency" | "usage" | "database" | "workers" | "cold-starts";

export interface AttentionItem {
  id: string;
  tone: "warning" | "destructive";
  /** One plain sentence: what is off, with the numbers. */
  text: string;
  tab: MetricsTab;
}

/** Cold starts failing more often than this, over at least COLD_START_MIN starts, are flagged. */
export const COLD_START_FAILURE_SHARE = 0.1;
const COLD_START_MIN = 5;
/** Workers erroring on more than this share of requests are flagged. */
export const WORKER_ERROR_SHARE = 0.01;
/** Usage is flagged when it is on course to pass this share of the plan by the end of the cycle. */
export const USAGE_PROJECTED_SHARE = 0.8;

const name = (id: ServiceId) => findService(id)?.name ?? id;

/** Share of the billing cycle gone by `now`, 0–1. */
export function cyclePace(cycle: MetricsPayload["cycle"], now = new Date()): number {
  const start = Date.parse(cycle.start);
  const end = Date.parse(cycle.end);
  return Math.min(1, Math.max(0, (now.getTime() - start) / (end - start)));
}

/**
 * Where usage lands at the end of the cycle if it keeps its pace so far, as a
 * share of what the plan includes. Early in a cycle there is too little to go on.
 */
export function projectedShare(usage: Usage, pace: number): number | null {
  if (pace < 0.1) return null;
  return usage.used / usage.included / pace;
}

/** Everything in the payload that deserves a look, worst first. */
export function needsAttention(data: MetricsPayload, now = new Date()): AttentionItem[] {
  const items: AttentionItem[] = [];

  for (const service of data.services) {
    if (service.status === "degraded") {
      items.push({
        id: `degraded-${service.id}`,
        tone: "destructive",
        text: `${name(service.id)} is degraded: ${service.errors24h ?? 0} server errors and a p95 of ${formatMs(service.p95Ms24h)} in the last 24 hours.`,
        tab: "services",
      });
    } else if ((service.errors24h ?? 0) > 0) {
      const errors = service.errors24h!;
      items.push({
        id: `errors-${service.id}`,
        tone: "warning",
        text: `${name(service.id)} answered ${errors} ${errors === 1 ? "request" : "requests"} with a server error in the last 24 hours.`,
        tab: "services",
      });
    }
  }

  for (const service of data.services) {
    const week = summarize(data.coldStarts, service.id, 7, now);
    if (week.count >= COLD_START_MIN && week.failed / week.count > COLD_START_FAILURE_SHARE) {
      items.push({
        id: `cold-${service.id}`,
        tone: week.failed / week.count > 0.5 ? "destructive" : "warning",
        text: `${name(service.id)} had ${week.failed} of ${week.count} cold starts fail in the last 7 days.`,
        tab: "cold-starts",
      });
    }
  }

  for (const worker of data.workers) {
    if (worker.kind === "other" || worker.requests === 0) continue;
    const share = worker.errors / worker.requests;
    if (share > WORKER_ERROR_SHARE) {
      items.push({
        id: `worker-${worker.script}`,
        tone: "warning",
        text: `The ${worker.script} Worker threw on ${formatPercent(share, 1)} of its requests (${formatCount(worker.errors)} of ${formatCount(worker.requests)}) this cycle.`,
        tab: "workers",
      });
    }
  }

  const pace = cyclePace(data.cycle, now);
  for (const usage of data.usage) {
    const projected = projectedShare(usage, pace);
    const limit = usage.stopsAt ? usage.stopsAt / usage.included : 1;
    if (projected !== null && projected > Math.min(USAGE_PROJECTED_SHARE, limit)) {
      items.push({
        id: `usage-${usage.resource}`,
        tone: projected > limit ? "destructive" : "warning",
        text: `${PLAN[usage.resource].label} is on course for ${formatPercent(projected, 0)} of the plan by the end of the cycle.`,
        tab: "usage",
      });
    }
  }

  for (const endpoint of data.endpoints) {
    if ((endpoint.p95Ms ?? 0) > DEGRADED_P95_MS) {
      items.push({
        id: `slow-${endpoint.service}-${endpoint.method}-${endpoint.route}`,
        tone: "warning",
        text: `${endpoint.method} ${endpoint.route} on ${name(endpoint.service)} takes ${formatMs(endpoint.p95Ms)} at p95.`,
        tab: "latency",
      });
    }
  }

  return items.sort((a, b) => (a.tone === b.tone ? 0 : a.tone === "destructive" ? -1 : 1));
}
