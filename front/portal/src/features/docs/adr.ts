import type { DocLink } from "./catalog";
import type { DocMeta } from "./types";

// Helpers for the ADR shape: dates, the areas a scope touches, and the
// links between decisions.

/** "2026-11-05" → "5 Nov 2026"; anything else comes back as given. */
export function formatDate(date: string | undefined, withYear = true): string | undefined {
  if (!date || !/^\d{4}-\d{2}-\d{2}$/.test(date)) return date;
  return new Date(`${date}T00:00:00Z`).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: withYear ? "numeric" : undefined,
    timeZone: "UTC",
  });
}

/** "2026-11-05" → "November 2026" */
export function monthOf(date: string | undefined): string {
  if (!date || !/^\d{4}-\d{2}/.test(date)) return "Undated";
  return new Date(`${date.slice(0, 7)}-01T00:00:00Z`).toLocaleDateString("en-GB", {
    month: "long",
    year: "numeric",
    timeZone: "UTC",
  });
}

export type Area = "backends" | "go-libs" | "ts-libs" | "sdks" | "fronts" | "tooling";

export const AREAS: readonly { id: Area; label: string }[] = [
  { id: "backends", label: "Backends" },
  { id: "go-libs", label: "Go libs" },
  { id: "ts-libs", label: "TS libs" },
  { id: "sdks", label: "SDKs" },
  { id: "fronts", label: "Fronts" },
  { id: "tooling", label: "Tooling" },
];

function areaOf(path: string): Area {
  if (path.startsWith("api/")) return "backends";
  if (path.startsWith("lib/go")) return "go-libs";
  if (path.startsWith("lib/ts")) return "ts-libs";
  if (path.startsWith("sdk/")) return "sdks";
  if (path === "front" || path.startsWith("front/")) return "fronts";
  // tools/, CI, hooks, e2e, root files and outside repos.
  return "tooling";
}

export function areasOf(meta: DocMeta): Set<Area> {
  return new Set((meta.scope ?? []).map(areaOf));
}

export type Relation = "supersedes" | "superseded_by" | "amends" | "amended_by";

export const RELATION_LABELS: Record<Relation, string> = {
  superseded_by: "Replaced by",
  supersedes: "Replaces",
  amended_by: "Amended by",
  amends: "Amends",
};

/** Every relation the ADR declares, in reading order, as [kind, numbers]. */
export function relationsOf(meta: DocMeta): [Relation, string[]][] {
  const out: [Relation, string[]][] = [];
  if (meta.superseded_by) out.push(["superseded_by", [meta.superseded_by]]);
  for (const kind of ["supersedes", "amended_by", "amends"] as const) {
    if (meta[kind]?.length) out.push([kind, meta[kind]!]);
  }
  return out;
}

export function isSuperseded(meta: DocMeta): boolean {
  return Boolean(meta.superseded_by) || meta.status?.toLowerCase().startsWith("superseded") === true;
}

/** Case-insensitive match on number, title and summary. */
export function matches(doc: DocLink, query: string): boolean {
  const needle = query.trim().toLowerCase();
  if (!needle) return true;
  const { number = "", title = "", summary = "" } = doc.meta;
  return `${number} ${title} ${summary}`.toLowerCase().includes(needle);
}

/** Groups docs (already in order) by the month they were decided. */
export function byMonth<T extends DocLink>(docs: readonly T[]): { month: string; docs: T[] }[] {
  const groups: { month: string; docs: T[] }[] = [];
  for (const doc of docs) {
    const month = monthOf(doc.meta.date);
    const last = groups.at(-1);
    if (last?.month === month) last.docs.push(doc);
    else groups.push({ month, docs: [doc] });
  }
  return groups;
}
