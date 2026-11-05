// Shared by the Worker (GitHub webhook, deploy hook, GET /api/delivery) and
// the Delivery page. Components are the part of a release tag before the
// version: `identityx/v0.42.0` → `identityx`, `pages/trieoh/v1.0.0` →
// `pages/trieoh`.

export type DeployStatus = "success" | "failure" | "cancelled";

export interface Deploy {
  component: string;
  tag: string;
  sha: string;
  actor: string;
  status: DeployStatus;
  at: string;
  runId: number | null;
}

export interface PullRequest {
  number: number;
  title: string;
  author: string;
  state: "open" | "closed" | "merged";
  draft: boolean;
  reviewState: string | null;
  openedAt: string;
  mergedAt: string | null;
}

export interface WorkflowRun {
  id: number;
  workflow: string;
  /** The branch, or the tag for release runs. */
  branch: string;
  sha: string;
  /** null while queued or running. */
  conclusion: string | null;
  attempt: number;
  queuedAt: string | null;
  startedAt: string | null;
  completedAt: string | null;
}

export interface Tag {
  name: string;
  component: string;
  version: string;
  sha: string;
  message: string | null;
  createdAt: string;
}

/** A release run, as deploy.yml's report job sent it to /webhooks/ci. */
export interface CiRun {
  runId: number;
  attempt: number;
  component: string;
  tag: string;
  /** service, front, sdk or pages */
  type: string;
  lane: string;
  status: DeployStatus;
  /** Where it stopped: parse, checks or publish. */
  stage: string;
  startedAt: string | null;
  completedAt: string | null;
}

/** How long one step of one pipeline takes, over its recent runs. */
export interface CiStepStat {
  type: string;
  /** checks or publish */
  job: string;
  /** The step's name; "Waiting for a runner" is the job's queue time. */
  step: string;
  runs: number;
  medianMs: number;
  p90Ms: number;
  lastMs: number;
  /** Up to the last 14 runs, oldest first. */
  recentMs: number[];
}

/** A number the actions recorded (.github/scripts/metric), over recent runs. */
export interface CiMetricStat {
  type: string;
  source: string;
  name: string;
  labels: Record<string, string>;
  runs: number;
  last: number;
  median: number;
}

export interface CiPayload {
  /** The 40 newest release runs, newest first. */
  runs: CiRun[];
  /** Per pipeline type, in step order. */
  steps: CiStepStat[];
  metrics: CiMetricStat[];
}

export interface DeliveryPayload {
  /** The latest deploy of each component. */
  production: Deploy[];
  /** Every deploy of the last 60 days, newest first, for the DORA numbers. */
  deploys: Deploy[];
  /** Open pull requests plus those merged in the last 30 days. */
  pulls: PullRequest[];
  /** The last 14 runs of each workflow, oldest first. */
  runs: WorkflowRun[];
  /** The 30 newest tags. */
  tags: Tag[];
  /** Release runs step by step, from /webhooks/ci. */
  ci: CiPayload;
}

const RELEASE_TAG = /^(.+)\/(v\d+\.\d+\.\d+.*)$/;

/** `univents-ui/v0.14.0` → { component: "univents-ui", version: "v0.14.0" } */
export function parseTag(name: string): { component: string; version: string } | null {
  const match = RELEASE_TAG.exec(name);
  return match ? { component: match[1], version: match[2] } : null;
}

export const REPO = "TrieOH/TheTree";
export const runUrl = (id: number) => `https://github.com/${REPO}/actions/runs/${id}`;
export const pullUrl = (number: number) => `https://github.com/${REPO}/pull/${number}`;
// TrieOH ships tags, not GitHub releases, so a tag links to its tree.
export const tagUrl = (name: string) => `https://github.com/${REPO}/tree/${name}`;
export const commitUrl = (sha: string) => `https://github.com/${REPO}/commit/${sha}`;
