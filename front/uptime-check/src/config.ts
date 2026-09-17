export interface Endpoint {
  /** Short slug used in KV keys, pings and the status page. */
  name: string;
  url: string;
}

/**
 * The public product APIs. This is the status-page checklist — the only
 * things we check and publish. Infra surfaces (forgejo, ntfy, staging…)
 * are deliberately not here.
 */
export const ENDPOINTS: Endpoint[] = [
  { name: "univents", url: "https://api.univents.com.br/health" },
  { name: "identityx", url: "https://api.identityx.com.br/health" },
  { name: "payssage", url: "https://api.payssage.trieoh.com/health" },
  { name: "informd", url: "https://api.informd.trieoh.com/health" },
];

export const CHECK = {
  /** Per-attempt fetch timeout. */
  timeoutMs: 10_000,
  /** Retries after the initial attempt fails (30s gaps, same invocation). */
  retries: 3,
  retryGapMs: 30_000,
} as const;

/**
 * Retention ladder. Raw 5-min points are kept 7 days, then compacted to
 * hourly averages (kept 30 days), then to daily averages (kept 1 year).
 * Degraded counts as up in all averages; buckets are colored 100% green,
 * 0% red, anything between amber.
 */
export const RETENTION = {
  rawDays: 7,
  hourlyDays: 30,
  dailyDays: 365,
} as const;

/** ntfy priorities for dead / recovery pings. Degraded is silent. */
export const NTFY = {
  deadPriority: "high",
  deadTags: ["rotating_light"],
  recoveryPriority: "default",
  recoveryTags: ["white_check_mark"],
} as const;

/** Must match wrangler.jsonc `triggers.crons`. */
export const CRON = {
  checks: "*/30 * * * *",
  compaction: "40 3 * * *",
} as const;
