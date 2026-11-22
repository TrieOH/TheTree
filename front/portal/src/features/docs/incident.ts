import type { DocMeta } from "./types";

// Helpers for the incident-report shape: severity, the phases an incident went
// through and how long each took. Times are UTC; a phase may be a bare date.

export type Severity = "high" | "medium" | "low";

export function severityOf(meta: DocMeta): Severity | undefined {
  const word = meta.severity?.toLowerCase();
  return word === "high" || word === "medium" || word === "low" ? word : undefined;
}

/** "medium" → "Medium severity" */
export function severityLabel(level: Severity): string {
  return `${level[0].toUpperCase()}${level.slice(1)} severity`;
}

export const SEVERITY_BADGE = { high: "destructive", medium: "warning", low: "muted" } as const;

const PHASES = [
  ["started", "Started"],
  ["detected", "Noticed"],
  ["mitigated", "Contained"],
  ["ended", "Over"],
] as const;

export interface Phase {
  label: string;
  /** As written in the preamble. */
  at: string;
  /** Whether `at` carries a time of day, not just a date. */
  timed: boolean;
  ms: number;
}

function parse(at: string): { ms: number; timed: boolean } | undefined {
  const timed = at.includes("T");
  const ms = Date.parse(timed ? at : `${at}T00:00:00Z`);
  return Number.isNaN(ms) ? undefined : { ms, timed };
}

/** The phases the report gives, in order. */
export function phasesOf(meta: DocMeta): Phase[] {
  return PHASES.flatMap(([key, label]) => {
    const at = meta[key];
    const parsed = at ? parse(at) : undefined;
    return parsed ? [{ label, at: at!, ...parsed }] : [];
  });
}

/** "26 Aug, 17:04 UTC" or "19 Oct" for a bare date; the year only when asked. */
export function formatMoment(phase: Pick<Phase, "ms" | "timed">, withYear = false): string {
  const date = new Date(phase.ms);
  const day = date.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: withYear ? "numeric" : undefined,
    timeZone: "UTC",
  });
  if (!phase.timed) return day;
  const time = date.toLocaleTimeString("en-GB", { hour: "2-digit", minute: "2-digit", timeZone: "UTC" });
  return `${day}, ${time} UTC`;
}

/** "13 min", "6 h 10 min", "14 days". Between a bare date and a time, only days are honest. */
export function formatSpan(from: Phase, to: Phase): string | undefined {
  const ms = to.ms - from.ms;
  if (ms < 0) return undefined;
  const minutes = Math.round(ms / 60_000);
  const days = Math.round(ms / 86_400_000);
  if (!from.timed || !to.timed) return days <= 1 ? "about a day" : `about ${days} days`;
  if (minutes < 60) return `${minutes} min`;
  if (minutes < 48 * 60) {
    const h = Math.floor(minutes / 60);
    const m = minutes % 60;
    return m ? `${h} h ${m} min` : `${h} h`;
  }
  return `${days} days`;
}

/** From the first phase to "Over", when both are known. */
export function durationOf(meta: DocMeta): string | undefined {
  const phases = phasesOf(meta);
  const first = phases[0];
  const last = phases.at(-1);
  return first && last && last.label === "Over" && first !== last ? formatSpan(first, last) : undefined;
}

const SERVICE_NAMES: Record<string, string> = {
  identityx: "IdentityX",
  payssage: "Payssage",
  univents: "Univents",
  informd: "Informd",
  forgejo: "Forgejo",
  "forgejo-runner": "Forgejo runner",
};

/** "identityx" → "IdentityX"; anything unknown stays as written. */
export function serviceName(id: string): string {
  return SERVICE_NAMES[id.toLowerCase()] ?? id;
}
