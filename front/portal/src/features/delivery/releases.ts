import { parseTag, type CiRun, type DeliveryPayload } from "./model";

// A merge to main is released as one tag per component it touches, and each
// tag runs deploy.yml on its own. Grouping the runs back by commit shows a
// release the way it was made: one change, the components it shipped, and how
// long the whole batch took from the first run starting to the last finishing.

export interface Release {
  /** The commit, or the tag when its commit isn't known. */
  key: string;
  sha: string | null;
  /** The pull request's title, else the commit or tag message, else the tag. */
  title: string;
  pr: number | null;
  /** Oldest start first. */
  runs: CiRun[];
  startedAt: string | null;
  completedAt: string | null;
}

export type ReleaseOutcome = "shipped" | "failed" | "cancelled" | "mixed";

const MERGE = /^Merge pull request #(\d+)/;

export function groupReleases(data: Pick<DeliveryPayload, "ci" | "tags" | "deploys" | "pulls">): Release[] {
  const shaOf = new Map<string, string>();
  for (const deploy of data.deploys) shaOf.set(deploy.tag, deploy.sha);
  for (const tag of data.tags) shaOf.set(tag.name, tag.sha);
  const messageOf = new Map(data.tags.filter((t) => t.message).map((t) => [t.sha, t.message!.split("\n")[0]]));
  const titleOf = new Map(data.pulls.map((pr) => [pr.number, pr.title]));

  const groups = new Map<string, CiRun[]>();
  for (const run of data.ci.runs) {
    const key = shaOf.get(run.tag) ?? run.tag;
    groups.set(key, [...(groups.get(key) ?? []), run]);
  }

  const time = (iso: string | null) => (iso ? Date.parse(iso) : NaN);
  return [...groups.entries()]
    .map(([key, runs]) => {
      const sha = shaOf.has(runs[0].tag) ? key : null;
      const message = sha ? messageOf.get(sha) : undefined;
      const pr = message ? Number(MERGE.exec(message)?.[1]) || null : null;
      const starts = runs.map((r) => time(r.startedAt)).filter(Number.isFinite);
      const ends = runs.map((r) => time(r.completedAt)).filter(Number.isFinite);
      return {
        key,
        sha,
        pr,
        title:
          (pr && titleOf.get(pr)) || (pr ? `Pull request #${pr}` : message) || runs.map((r) => r.tag).join(", "),
        runs: [...runs].sort((a, b) => time(a.startedAt) - time(b.startedAt)),
        startedAt: starts.length ? new Date(Math.min(...starts)).toISOString() : null,
        completedAt: ends.length ? new Date(Math.max(...ends)).toISOString() : null,
      };
    })
    .sort((a, b) => time(b.startedAt) - time(a.startedAt));
}

export function releaseOutcome(release: Release): ReleaseOutcome {
  const statuses = new Set(release.runs.map((r) => r.status));
  if (statuses.size > 1) return "mixed";
  if (statuses.has("failure")) return "failed";
  if (statuses.has("cancelled")) return "cancelled";
  return "shipped";
}

/** Where one run sits on its release's clock, as shares 0–1 of the release's span. */
export function runSpan(release: Release, run: CiRun): { start: number; width: number } | null {
  if (!release.startedAt || !release.completedAt || !run.startedAt || !run.completedAt) return null;
  const from = Date.parse(release.startedAt);
  const span = Date.parse(release.completedAt) - from;
  if (span <= 0) return { start: 0, width: 1 };
  return {
    start: (Date.parse(run.startedAt) - from) / span,
    width: (Date.parse(run.completedAt) - Date.parse(run.startedAt)) / span,
  };
}

/** `identityx/v0.49.1` → `v0.49.1`, for a run already labelled with its component. */
export const versionOf = (tag: string) => parseTag(tag)?.version ?? tag;
