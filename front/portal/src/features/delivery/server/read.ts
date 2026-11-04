import type { DeliveryPayload, Deploy, DeployStatus, PullRequest, Tag, WorkflowRun } from "../model";

// GET /api/delivery: one small payload for the Delivery page and the service
// pages' current tag.

interface DeployRow {
  component: string;
  tag: string;
  sha: string;
  actor: string;
  status: DeployStatus;
  health_ok: number | null;
  at: string;
  run_id: number | null;
}

interface PullRow {
  number: number;
  title: string;
  author: string;
  state: PullRequest["state"];
  draft: number;
  review_state: string | null;
  opened_at: string;
  merged_at: string | null;
}

interface RunRow {
  id: number;
  workflow: string;
  branch: string;
  sha: string;
  conclusion: string | null;
  attempt: number;
  queued_at: string | null;
  started_at: string | null;
  completed_at: string | null;
}

interface TagRow {
  name: string;
  component: string;
  version: string;
  sha: string;
  message: string | null;
  created_at: string;
}

const DEPLOY_COLUMNS = "component, tag, sha, actor, status, health_ok, at, run_id";

function toDeploy(row: DeployRow): Deploy {
  return {
    component: row.component,
    tag: row.tag,
    sha: row.sha,
    actor: row.actor,
    status: row.status,
    healthOk: row.health_ok === null ? null : row.health_ok === 1,
    at: row.at,
    runId: row.run_id,
  };
}

export async function readDelivery(db: D1Database, now = new Date()): Promise<DeliveryPayload> {
  const days = (n: number) => new Date(now.getTime() - n * 86_400_000).toISOString();

  const [production, deploys, pulls, runs, tags] = await db.batch<unknown>([
    db.prepare(
      `SELECT ${DEPLOY_COLUMNS} FROM (
         SELECT *, ROW_NUMBER() OVER (PARTITION BY component ORDER BY at DESC, id DESC) AS n FROM deploys
       ) WHERE n = 1 ORDER BY component`,
    ),
    db.prepare(`SELECT ${DEPLOY_COLUMNS} FROM deploys WHERE at >= ? ORDER BY at DESC, id DESC`).bind(days(60)),
    db.prepare(
      `SELECT number, title, author, state, draft, review_state, opened_at, merged_at FROM pull_requests
       WHERE state = 'open' OR (state = 'merged' AND merged_at >= ?)
       ORDER BY state = 'open' DESC, COALESCE(merged_at, opened_at) DESC`,
    ).bind(days(30)),
    db.prepare(
      `SELECT id, workflow, branch, sha, conclusion, attempt, queued_at, started_at, completed_at FROM (
         SELECT *, ROW_NUMBER() OVER (PARTITION BY workflow ORDER BY queued_at DESC) AS n FROM workflow_runs
       ) WHERE n <= 14 ORDER BY workflow, queued_at`,
    ),
    db.prepare(
      "SELECT name, component, version, sha, message, created_at FROM tags ORDER BY created_at DESC LIMIT 30",
    ),
  ]);

  return {
    production: (production.results as DeployRow[]).map(toDeploy),
    deploys: (deploys.results as DeployRow[]).map(toDeploy),
    pulls: (pulls.results as PullRow[]).map(
      (row): PullRequest => ({
        number: row.number,
        title: row.title,
        author: row.author,
        state: row.state,
        draft: row.draft === 1,
        reviewState: row.review_state,
        openedAt: row.opened_at,
        mergedAt: row.merged_at,
      }),
    ),
    runs: (runs.results as RunRow[]).map(
      (row): WorkflowRun => ({
        id: row.id,
        workflow: row.workflow,
        branch: row.branch,
        sha: row.sha,
        conclusion: row.conclusion,
        attempt: row.attempt,
        queuedAt: row.queued_at,
        startedAt: row.started_at,
        completedAt: row.completed_at,
      }),
    ),
    tags: (tags.results as TagRow[]).map(
      (row): Tag => ({
        name: row.name,
        component: row.component,
        version: row.version,
        sha: row.sha,
        message: row.message,
        createdAt: row.created_at,
      }),
    ),
  };
}
