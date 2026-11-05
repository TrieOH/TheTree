import type { CiMetricStat, CiPayload, CiRun, CiStepStat, DeployStatus } from "../model";

// The CI part of GET /api/delivery: the newest release runs, and per pipeline
// type how long each step takes and what the actions measured, over the last
// RECENT runs of that type.

const RECENT = 20;
const SPARK = 14;
export const QUEUE_STEP = "Waiting for a runner";

interface RunRow {
  run_id: number;
  attempt: number;
  component: string;
  tag: string;
  type: string;
  lane: string;
  status: DeployStatus;
  stage: string;
  started_at: string | null;
  completed_at: string | null;
}

interface StepRow {
  type: string;
  run_at: string | null;
  job: string;
  number: number;
  name: string;
  status: string | null;
  started_at: string | null;
  completed_at: string | null;
}

interface MetricRow {
  type: string;
  run_at: string | null;
  source: string;
  name: string;
  value: number;
  labels: string;
}

// recent runs per type, oldest first within each type
const RECENT_RUNS = `
  WITH recent AS (
    SELECT run_id, attempt, type, started_at AS run_at FROM (
      SELECT *, ROW_NUMBER() OVER (PARTITION BY type ORDER BY started_at DESC) AS n FROM ci_runs
    ) WHERE n <= ${RECENT}
  )`;

function percentile(values: number[], p: number): number {
  const sorted = [...values].sort((a, b) => a - b);
  return sorted[Math.min(sorted.length - 1, Math.max(0, Math.ceil(p * sorted.length) - 1))];
}

const span = (from: string | null, to: string | null) =>
  from && to ? Math.max(0, Date.parse(to) - Date.parse(from)) : null;

/** "front / checks" → "checks" */
const jobOf = (name: string) => name.split(" / ").at(-1)!;

const JOB_ORDER = ["parse", "checks", "publish", "report"];

export async function readCi(db: D1Database): Promise<CiPayload> {
  const [runs, steps, jobs, metrics] = await db.batch<unknown>([
    db.prepare(
      `SELECT run_id, attempt, component, tag, type, lane, status, stage, started_at, completed_at
       FROM ci_runs ORDER BY started_at DESC LIMIT 40`,
    ),
    db.prepare(
      `${RECENT_RUNS}
       SELECT r.type, r.run_at, s.job, s.number, s.name, s.status, s.started_at, s.completed_at
       FROM ci_steps s JOIN recent r USING (run_id, attempt)
       ORDER BY r.run_at, s.job, s.number`,
    ),
    // a job's queue time, as a step before its first one
    db.prepare(
      `${RECENT_RUNS}
       SELECT r.type, r.run_at, j.name AS job, -1 AS number, '${QUEUE_STEP}' AS name, j.status,
              j.queued_at AS started_at, j.started_at AS completed_at
       FROM ci_jobs j JOIN recent r USING (run_id, attempt)
       ORDER BY r.run_at`,
    ),
    db.prepare(
      `${RECENT_RUNS}
       SELECT r.type, r.run_at, m.source, m.name, m.value, m.labels
       FROM ci_metrics m JOIN recent r USING (run_id, attempt)
       ORDER BY r.run_at`,
    ),
  ]);

  return {
    runs: (runs.results as RunRow[]).map(
      (row): CiRun => ({
        runId: row.run_id,
        attempt: row.attempt,
        component: row.component,
        tag: row.tag,
        type: row.type,
        lane: row.lane,
        status: row.status,
        stage: row.stage,
        startedAt: row.started_at,
        completedAt: row.completed_at,
      }),
    ),
    steps: stepStats([...(jobs.results as StepRow[]), ...(steps.results as StepRow[])]),
    metrics: metricStats(metrics.results as MetricRow[]),
  };
}

/** Rows oldest run first; skipped steps and steps without times are left out. */
export function stepStats(rows: StepRow[]): CiStepStat[] {
  const groups = new Map<string, { type: string; job: string; step: string; order: number[]; ms: number[] }>();
  for (const row of [...rows].sort((a, b) => (a.run_at ?? "").localeCompare(b.run_at ?? ""))) {
    const ms = span(row.started_at, row.completed_at);
    if (ms === null || row.status === "skipped") continue;
    const job = jobOf(row.job);
    const key = `${row.type}\u0000${job}\u0000${row.name}`;
    const group = groups.get(key) ?? { type: row.type, job, step: row.name, order: [], ms: [] };
    group.order.push(row.number);
    group.ms.push(ms);
    groups.set(key, group);
  }
  return [...groups.values()]
    .map((group) => ({
      stat: {
        type: group.type,
        job: group.job,
        step: group.step,
        runs: group.ms.length,
        medianMs: percentile(group.ms, 0.5),
        p90Ms: percentile(group.ms, 0.9),
        lastMs: group.ms.at(-1)!,
        recentMs: group.ms.slice(-SPARK),
      },
      order: percentile(group.order, 0.5),
    }))
    .sort(
      (a, b) =>
        a.stat.type.localeCompare(b.stat.type) ||
        JOB_ORDER.indexOf(a.stat.job) - JOB_ORDER.indexOf(b.stat.job) ||
        a.order - b.order,
    )
    .map(({ stat }) => stat);
}

/** Rows oldest run first. */
export function metricStats(rows: MetricRow[]): CiMetricStat[] {
  const groups = new Map<string, { type: string; source: string; name: string; labels: string; values: number[] }>();
  for (const row of rows) {
    const key = `${row.type}\u0000${row.source}\u0000${row.name}\u0000${row.labels}`;
    const group = groups.get(key) ?? { type: row.type, source: row.source, name: row.name, labels: row.labels, values: [] };
    group.values.push(row.value);
    groups.set(key, group);
  }
  return [...groups.values()]
    .map(
      (group): CiMetricStat => ({
        type: group.type,
        source: group.source,
        name: group.name,
        labels: JSON.parse(group.labels) as Record<string, string>,
        runs: group.values.length,
        last: group.values.at(-1)!,
        median: percentile(group.values, 0.5),
      }),
    )
    .sort(
      (a, b) =>
        a.type.localeCompare(b.type) ||
        a.source.localeCompare(b.source) ||
        a.name.localeCompare(b.name) ||
        JSON.stringify(a.labels).localeCompare(JSON.stringify(b.labels)),
    );
}
