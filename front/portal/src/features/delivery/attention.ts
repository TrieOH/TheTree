import type { AttentionItem } from "@/shared/ui/attention";
import { parseTag, type DeliveryPayload } from "./model";
import { isSlowLastRun, stepLabel } from "./pipeline";
import { formatStepTime, timeAgo } from "./stats";

/** Which Delivery tab explains an item. */
export type DeliveryTab = "production" | "releases" | "pipelines" | "pulls";

const HOUR = 3_600_000;
const DAY = 24 * HOUR;
/** An open pull request waiting this long for a review is flagged. */
export const REVIEW_WAIT_MS = 2 * DAY;
/** A tag with no release run this long after it was pushed is flagged, up to TAG_LOOKBACK_MS. */
const TAG_GRACE_MS = HOUR;
const TAG_LOOKBACK_MS = 7 * DAY;

const PIPELINE: Record<string, string> = { service: "service", front: "front", sdk: "SDK", pages: "Pages" };

const days = (ms: number) => {
  const n = Math.floor(ms / DAY);
  return `${n} ${n === 1 ? "day" : "days"}`;
};

/** Everything in the payload that deserves a look, worst first. */
export function needsAttention(data: DeliveryPayload, now = new Date()): AttentionItem<DeliveryTab>[] {
  const items: AttentionItem<DeliveryTab>[] = [];
  const at = now.getTime();

  // A component whose newest release run didn't ship. Runs are newest first.
  const seen = new Set<string>();
  for (const run of data.ci.runs) {
    if (seen.has(run.component)) continue;
    seen.add(run.component);
    if (run.status !== "failure") continue;
    const version = parseTag(run.tag)?.version ?? run.tag;
    const where = run.stage === "publish" ? "failed while publishing" : `was blocked at ${run.stage}`;
    const when = run.startedAt ? ` ${timeAgo(run.startedAt, now)}` : "";
    items.push({
      id: `release-${run.component}`,
      tone: "destructive",
      text: `${run.component} ${version} ${where}${when}, and nothing newer has shipped.`,
      tab: "releases",
    });
  }
  for (const deploy of data.production) {
    if (deploy.status !== "failure" || seen.has(deploy.component)) continue;
    items.push({
      id: `deploy-${deploy.component}`,
      tone: "destructive",
      text: `The last deploy of ${deploy.component} failed ${timeAgo(deploy.at, now)}.`,
      tab: "production",
    });
  }

  // A tag that never reported a run: the workflow didn't start or didn't call back.
  // Only tags newer than the oldest run we hold can be judged.
  const reported = new Set([...data.ci.runs.map((r) => r.tag), ...data.deploys.map((d) => d.tag)]);
  const starts = data.ci.runs.flatMap((r) => (r.startedAt ? [Date.parse(r.startedAt)] : []));
  const horizon = Math.max(at - TAG_LOOKBACK_MS, starts.length ? Math.min(...starts) : at);
  for (const tag of data.tags) {
    const created = Date.parse(tag.createdAt);
    if (reported.has(tag.name) || created < horizon || at - created < TAG_GRACE_MS) continue;
    items.push({
      id: `tag-${tag.name}`,
      tone: "warning",
      text: `${tag.name} was tagged ${timeAgo(tag.createdAt, now)} but no release run reported back.`,
      tab: "releases",
    });
  }

  // Security findings in the newest scan of each pipeline.
  for (const type of Object.keys(PIPELINE)) {
    const latest = (severity: string) =>
      data.ci.metrics
        .filter((m) => m.type === type && m.name === "findings" && m.labels.severity === severity)
        .reduce((sum, m) => sum + m.last, 0);
    const critical = latest("CRITICAL");
    const high = latest("HIGH");
    if (critical + high === 0) continue;
    const parts = [critical && `${critical} critical`, high && `${high} high`].filter(Boolean).join(" and ");
    items.push({
      id: `trivy-${type}`,
      tone: critical ? "destructive" : "warning",
      text: `The last ${PIPELINE[type]} release scan found ${parts} severity ${critical + high === 1 ? "vulnerability" : "vulnerabilities"}.`,
      tab: "pipelines",
    });
  }

  // Steps whose last run was far slower than usual.
  for (const step of data.ci.steps) {
    if (!PIPELINE[step.type] || !isSlowLastRun(step)) continue;
    items.push({
      id: `slow-${step.type}-${step.job}-${step.step}`,
      tone: "warning",
      text: `${stepLabel(step.step)} took ${formatStepTime(step.lastMs)} in the last ${PIPELINE[step.type]} release, against a usual ${formatStepTime(step.medianMs)}.`,
      tab: "pipelines",
    });
  }

  // Pull requests left waiting for a review.
  for (const pr of data.pulls) {
    if (pr.state !== "open" || pr.draft || pr.reviewState === "approved" || pr.reviewState === "changes_requested") continue;
    const waited = at - Date.parse(pr.openedAt);
    if (waited < REVIEW_WAIT_MS) continue;
    items.push({
      id: `pr-${pr.number}`,
      tone: "warning",
      text: `#${pr.number} “${pr.title}” has waited ${days(waited)} for a review.`,
      tab: "pulls",
    });
  }

  return items.sort((a, b) => (a.tone === b.tone ? 0 : a.tone === "destructive" ? -1 : 1));
}
