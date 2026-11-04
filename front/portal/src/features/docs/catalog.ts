import type { Component } from "solid-js";
import type { DocMeta } from "./types";

// Docs come straight from the repo's docs/ folder at build time; nothing is
// copied into the portal. Decisions are the ADRs; the other categories fill
// in as reports and runbooks are converted (phase 4).

export type DocCategoryId = "decisions" | "runbooks" | "incident-reports" | "onboarding";

export interface DocCategory {
  id: DocCategoryId;
  label: string;
  description: string;
}

export const DOC_CATEGORIES: readonly DocCategory[] = [
  { id: "decisions", label: "Decisions", description: "Architecture decision records: why things are this way." },
  { id: "runbooks", label: "Runbooks", description: "How to deploy, roll back and fix the usual things." },
  { id: "incident-reports", label: "Incident reports", description: "What broke, why, and what changed after." },
  { id: "onboarding", label: "Onboarding", description: "Run TheTree locally and learn how we ship." },
];

type Loader = () => Promise<{ default: Component<{ components?: Record<string, unknown> }> }>;

export interface DocEntry {
  category: DocCategoryId;
  slug: string;
  meta: DocMeta;
  load: Loader;
}

const adrMeta = import.meta.glob<DocMeta>("../../../../../docs/adr/*.md", {
  eager: true,
  import: "meta",
});
const adrContent = import.meta.glob<Awaited<ReturnType<Loader>>>("../../../../../docs/adr/*.md");

function slugOf(path: string): string {
  return path.slice(path.lastIndexOf("/") + 1).replace(/\.mdx?$/, "").toLowerCase();
}

const decisions: DocEntry[] = Object.entries(adrMeta)
  .map(([path, meta]) => ({
    category: "decisions" as const,
    slug: slugOf(path),
    meta,
    load: adrContent[path] as Loader,
  }))
  // Newest decision first.
  .sort((a, b) => b.slug.localeCompare(a.slug));

const DOCS: Record<DocCategoryId, readonly DocEntry[]> = {
  decisions,
  runbooks: [],
  "incident-reports": [],
  onboarding: [],
};

export function findCategory(id: string): DocCategory | undefined {
  return DOC_CATEGORIES.find((category) => category.id === id);
}

export function docsIn(category: DocCategoryId): readonly DocEntry[] {
  return DOCS[category];
}

export function findDoc(category: string, slug: string): DocEntry | undefined {
  const known = findCategory(category);
  return known ? DOCS[known.id].find((doc) => doc.slug === slug) : undefined;
}

/** "accepted; `repo/` amended …" → "accepted" */
export function statusWord(status: string | undefined): string | undefined {
  return status?.match(/^[a-z]+/i)?.[0].toLowerCase();
}
