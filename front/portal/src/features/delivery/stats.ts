import type { Deploy, PullRequest, WorkflowRun } from "./model";

// The four DORA numbers over the payload, plus CI durations.

const DAY = 86_400_000;

function median(values: number[]): number | null {
  if (values.length === 0) return null;
  const sorted = [...values].sort((a, b) => a - b);
  const mid = Math.floor(sorted.length / 2);
  return sorted.length % 2 ? sorted[mid] : (sorted[mid - 1] + sorted[mid]) / 2;
}

/** A deploy that went wrong: its workflow run failed. */
export function isFailed(deploy: Deploy): boolean {
  return deploy.status === "failure";
}

const within = (iso: string, days: number, now: Date) => Date.parse(iso) >= now.getTime() - days * DAY;

/** Successful deploys per week over the last 28 days. */
export function deploysPerWeek(deploys: readonly Deploy[], now = new Date()): number {
  return deploys.filter((d) => d.status === "success" && within(d.at, 28, now)).length / 4;
}

/** Median ms from opening a PR to merging it, over the last 30 days. */
export function prLeadTime(pulls: readonly PullRequest[], now = new Date()): number | null {
  return median(
    pulls
      .filter((pr) => pr.mergedAt && within(pr.mergedAt, 30, now))
      .map((pr) => Date.parse(pr.mergedAt!) - Date.parse(pr.openedAt)),
  );
}

/** Share of the last 30 days' deploys that failed; null without deploys. */
export function changeFailureRate(deploys: readonly Deploy[], now = new Date()): number | null {
  const recent = deploys.filter((d) => d.status !== "cancelled" && within(d.at, 30, now));
  return recent.length ? recent.filter(isFailed).length / recent.length : null;
}

/**
 * Median ms from a failed deploy to the next good one of the same component,
 * over the last 30 days. Failures not yet fixed are left out.
 */
export function timeToRestore(deploys: readonly Deploy[], now = new Date()): number | null {
  const oldestFirst = [...deploys].sort((a, b) => Date.parse(a.at) - Date.parse(b.at));
  const restores: number[] = [];
  for (const [i, deploy] of oldestFirst.entries()) {
    if (!isFailed(deploy) || !within(deploy.at, 30, now)) continue;
    const previous = oldestFirst.slice(0, i).findLast((d) => d.component === deploy.component);
    if (previous && isFailed(previous)) continue; // counted from the first failure in a row
    const fix = oldestFirst
      .slice(i + 1)
      .find((d) => d.component === deploy.component && d.status === "success" && !isFailed(d));
    if (fix) restores.push(Date.parse(fix.at) - Date.parse(deploy.at));
  }
  return median(restores);
}

/** Median ms a workflow's completed runs took, from start to finish. */
export function typicalDuration(runs: readonly WorkflowRun[]): number | null {
  return median(
    runs
      .filter((run) => run.startedAt && run.completedAt)
      .map((run) => Date.parse(run.completedAt!) - Date.parse(run.startedAt!)),
  );
}

export function formatDuration(ms: number | null): string {
  if (ms === null) return "—";
  const minutes = ms / 60_000;
  if (minutes < 1) return `${Math.round(ms / 1000)} s`;
  if (minutes < 60) return `${Math.round(minutes)} min`;
  const hours = minutes / 60;
  return hours < 48 ? `${hours.toFixed(1)} h` : `${(hours / 24).toFixed(1)} d`;
}

export function timeAgo(iso: string, now = new Date()): string {
  const minutes = Math.round((now.getTime() - Date.parse(iso)) / 60_000);
  if (minutes < 1) return "just now";
  if (minutes < 60) return `${minutes} min ago`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours} h ago`;
  return `${Math.round(hours / 24)} d ago`;
}
