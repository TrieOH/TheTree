import type { Component } from "solid-js";
import type { ServiceId } from "@/features/catalog/services";
import { docsIn, type DocLink } from "./catalog";
import type { DocMeta } from "./types";

// A service's own docs: docs/services/<service>/{consumer,maintainer}/*.md(x).
// Consumer pages are for whoever builds on the service, maintainer pages for
// whoever changes it. Inside each, a doc's `**Group:**` field puts it under a
// heading of the service's choosing, and `**Order:**` sets the reading order.

export type Audience = "consumer" | "maintainer";

export interface AudienceInfo {
  id: Audience;
  label: string;
  description: string;
}

export const AUDIENCES: readonly AudienceInfo[] = [
  { id: "consumer", label: "Consumer", description: "For apps and services that call it." },
  { id: "maintainer", label: "Maintainer", description: "For whoever changes the service itself." },
];

type Loader = () => Promise<{ default: Component<{ components?: Record<string, unknown> }> }>;

export interface ServiceDoc {
  service: string;
  audience: Audience;
  slug: string;
  meta: DocMeta;
  load: Loader;
}

/** A service doc without its loader: safe to hand through route loader data. */
export type ServiceDocLink = Omit<ServiceDoc, "load">;

const meta = import.meta.glob<DocMeta>("../../../../../docs/services/*/{consumer,maintainer}/*.{md,mdx}", {
  eager: true,
  import: "meta",
});
const content = import.meta.glob<Awaited<ReturnType<Loader>>>(
  "../../../../../docs/services/*/{consumer,maintainer}/*.{md,mdx}",
);

function rank(doc: ServiceDoc): number {
  const order = Number(doc.meta.order);
  return Number.isFinite(order) && doc.meta.order ? order : Infinity;
}

const AUDIENCE_RANK: Record<Audience, number> = { consumer: 0, maintainer: 1 };

// Consumer pages first, then maintainer pages, each in reading order.
const ALL: ServiceDoc[] = Object.entries(meta)
  .map(([path, docMeta]) => {
    const [service, audience, file] = path.split("/").slice(-3);
    return {
      service,
      audience: audience as Audience,
      slug: file.replace(/\.mdx?$/, "").toLowerCase(),
      meta: docMeta,
      load: content[path] as Loader,
    };
  })
  .sort(
    (a, b) =>
      a.service.localeCompare(b.service) ||
      AUDIENCE_RANK[a.audience] - AUDIENCE_RANK[b.audience] ||
      rank(a) - rank(b) ||
      (a.meta.title ?? a.slug).localeCompare(b.meta.title ?? b.slug),
  );

export function serviceDocs(service: ServiceId): readonly ServiceDoc[] {
  return ALL.filter((doc) => doc.service === service);
}

export function findServiceDoc(service: string, slug: string): ServiceDoc | undefined {
  return ALL.find((doc) => doc.service === service && doc.slug === slug);
}

export interface DocGroup {
  /** The doc's `**Group:**`, or "" for docs without one. */
  name: string;
  docs: ServiceDocLink[];
}

const link = (doc: ServiceDoc): ServiceDocLink => ({
  service: doc.service,
  audience: doc.audience,
  slug: doc.slug,
  meta: doc.meta,
});

/** One audience's docs under their groups, groups in the order their first doc is read. */
export function groupsOf(service: ServiceId, audience: Audience): DocGroup[] {
  const groups: DocGroup[] = [];
  for (const doc of serviceDocs(service)) {
    if (doc.audience !== audience) continue;
    const name = doc.meta.group ?? "";
    const group = groups.find((entry) => entry.name === name);
    if (group) group.docs.push(link(doc));
    else groups.push({ name, docs: [link(doc)] });
  }
  return groups;
}

/** The docs on either side of one, in the order the service's index lists them. */
export function serviceNeighbours(doc: ServiceDoc): { previous?: ServiceDocLink; next?: ServiceDocLink } {
  const service = doc.service as ServiceId;
  const ordered = AUDIENCES.flatMap((audience) => groupsOf(service, audience.id).flatMap((group) => group.docs));
  const at = ordered.findIndex((other) => other.slug === doc.slug);
  return at < 0 ? {} : { previous: ordered[at - 1], next: ordered[at + 1] };
}

/** ADRs whose scope names the service, like `api/identityx` or `api/{payssage,univents}`. */
export function decisionsFor(service: ServiceId): DocLink[] {
  return docsIn("decisions").filter((doc) => doc.meta.scope?.some((path) => path.toLowerCase().includes(service)));
}

/** Incident reports that list the service as affected. */
export function incidentsFor(service: ServiceId): DocLink[] {
  return docsIn("incident-reports").filter((doc) => doc.meta.services?.includes(service));
}
