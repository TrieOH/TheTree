import { timingSafeEqual } from "@/shared/lib/hmac";

// POST /webhooks/ci, sent by deploy.yml's report job (.github/scripts/ci-report.py)
// with `Authorization: Bearer $DEPLOY_HOOK_SECRET` after every release run:
// the run, each job and step with its start and end (from the Actions API),
// and the metrics the actions recorded. A resent attempt replaces itself.

export const CI_HOOK_PATH = "/webhooks/ci";

const STATUSES = new Set(["success", "failure", "cancelled"]);
const STAGES = new Set(["parse", "checks", "publish"]);
// far above a real run (~5 jobs, ~60 steps, ~80 metrics)
const MAX_JOBS = 30;
const MAX_STEPS = 100;
const MAX_SUBSTEPS = 300;
const MAX_METRICS = 1000;
// D1 binds at most 100 parameters per statement
const MAX_PARAMS = 100;

interface StepReport {
  number: number;
  name: string;
  status: string | null;
  started_at: string | null;
  completed_at: string | null;
}

/** An action inside a composite step, from the job's log. */
interface SubstepReport {
  step: number | null;
  action_id: string;
  name: string;
  depth: number;
  completed_at: string;
  duration_ms: number;
}

interface JobReport {
  name: string;
  status: string | null;
  queued_at: string | null;
  started_at: string | null;
  completed_at: string | null;
  runner: string | null;
  steps: StepReport[];
  /** Optional: reports from before the substeps went without. */
  substeps?: SubstepReport[];
}

interface MetricReport {
  job: string;
  source: string;
  name: string;
  value: number;
  labels: Record<string, string>;
}

interface CiReport {
  run_id: number;
  attempt: number;
  component: string;
  tag: string;
  sha: string;
  type: string;
  lane: string;
  status: string;
  stage: string;
  jobs: JobReport[];
  metrics: MetricReport[];
}

const str = (value: unknown) => typeof value === "string";
const optStr = (value: unknown) => value === null || value === undefined || typeof value === "string";

function isStep(value: unknown): value is StepReport {
  const step = value as Partial<StepReport> | null;
  return Number.isInteger(step?.number) && str(step!.name) && optStr(step!.status) && optStr(step!.started_at) && optStr(step!.completed_at);
}

function isSubstep(value: unknown): value is SubstepReport {
  const sub = value as Partial<SubstepReport> | null;
  return (
    (sub?.step === null || Number.isInteger(sub?.step)) &&
    str(sub!.action_id) &&
    str(sub!.name) &&
    Number.isInteger(sub!.depth) &&
    str(sub!.completed_at) &&
    Number.isInteger(sub!.duration_ms)
  );
}

function isJob(value: unknown): value is JobReport {
  const job = value as Partial<JobReport> | null;
  return (
    str(job?.name) &&
    optStr(job!.status) &&
    optStr(job!.queued_at) &&
    optStr(job!.started_at) &&
    optStr(job!.completed_at) &&
    optStr(job!.runner) &&
    Array.isArray(job!.steps) &&
    job!.steps.length <= MAX_STEPS &&
    job!.steps.every(isStep) &&
    (job!.substeps === undefined ||
      (Array.isArray(job!.substeps) && job!.substeps.length <= MAX_SUBSTEPS && job!.substeps.every(isSubstep)))
  );
}

function isMetric(value: unknown): value is MetricReport {
  const metric = value as Partial<MetricReport> | null;
  return (
    str(metric?.job) &&
    str(metric!.source) &&
    str(metric!.name) &&
    Number.isFinite(metric!.value) &&
    typeof metric!.labels === "object" &&
    metric!.labels !== null &&
    Object.values(metric!.labels).every(str)
  );
}

function isReport(value: unknown): value is CiReport {
  const report = value as Partial<CiReport> | null;
  return (
    Number.isInteger(report?.run_id) &&
    Number.isInteger(report!.attempt) &&
    [report!.component, report!.tag, report!.sha, report!.type, report!.lane].every(str) &&
    STATUSES.has(report!.status as string) &&
    STAGES.has(report!.stage as string) &&
    Array.isArray(report!.jobs) &&
    report!.jobs.length <= MAX_JOBS &&
    report!.jobs.every(isJob) &&
    Array.isArray(report!.metrics) &&
    report!.metrics.length <= MAX_METRICS &&
    report!.metrics.every(isMetric)
  );
}

/** INSERTs of as many rows per statement as D1's parameter limit allows. */
function inserts(db: D1Database, table: string, columns: string[], rows: unknown[][]): D1PreparedStatement[] {
  const perStatement = Math.floor(MAX_PARAMS / columns.length);
  const row = `(${columns.map(() => "?").join(", ")})`;
  const statements: D1PreparedStatement[] = [];
  for (let i = 0; i < rows.length; i += perStatement) {
    const chunk = rows.slice(i, i + perStatement);
    statements.push(
      db
        .prepare(`INSERT INTO ${table} (${columns.join(", ")}) VALUES ${chunk.map(() => row).join(", ")}`)
        .bind(...chunk.flat()),
    );
  }
  return statements;
}

/** The runner can reuse an action id (post steps); suffix repeats so each row keeps its own key. */
function uniqueIds(substeps: SubstepReport[]): SubstepReport[] {
  const seen = new Map<string, number>();
  return substeps.map((sub) => {
    const n = (seen.get(sub.action_id) ?? 0) + 1;
    seen.set(sub.action_id, n);
    return n === 1 ? sub : { ...sub, action_id: `${sub.action_id}#${n}` };
  });
}

const sortedLabels =(labels: Record<string, string>) =>
  JSON.stringify(Object.fromEntries(Object.entries(labels).sort(([a], [b]) => a.localeCompare(b))));

function bounds(jobs: JobReport[]): [string | null, string | null] {
  const starts = jobs.map((job) => job.started_at).filter((at): at is string => !!at).sort();
  const ends = jobs.map((job) => job.completed_at).filter((at): at is string => !!at).sort();
  return [starts[0] ?? null, ends.at(-1) ?? null];
}

export async function handleCiHook(request: Request, env: Env): Promise<Response> {
  const auth = request.headers.get("authorization") ?? "";
  if (!timingSafeEqual(auth, `Bearer ${env.DEPLOY_HOOK_SECRET}`)) {
    return new Response("unauthorized", { status: 401 });
  }

  let report: unknown;
  try {
    report = await request.json();
  } catch {
    return new Response("invalid body", { status: 400 });
  }
  if (!isReport(report)) return new Response("invalid report", { status: 400 });

  const db = env.DB;
  const key = [report.run_id, report.attempt];
  const [startedAt, completedAt] = bounds(report.jobs);

  await db.batch([
    ...["ci_runs", "ci_jobs", "ci_steps", "ci_substeps", "ci_metrics"].map((table) =>
      db.prepare(`DELETE FROM ${table} WHERE run_id = ? AND attempt = ?`).bind(...key),
    ),
    db
      .prepare(
        `INSERT INTO ci_runs (run_id, attempt, component, tag, sha, type, lane, status, stage, started_at, completed_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(
        ...key,
        report.component,
        report.tag,
        report.sha,
        report.type,
        report.lane,
        report.status,
        report.stage,
        startedAt,
        completedAt,
      ),
    ...inserts(
      db,
      "ci_jobs",
      ["run_id", "attempt", "name", "status", "queued_at", "started_at", "completed_at", "runner"],
      report.jobs.map((job) => [...key, job.name, job.status, job.queued_at, job.started_at, job.completed_at, job.runner]),
    ),
    ...inserts(
      db,
      "ci_steps",
      ["run_id", "attempt", "job", "number", "name", "status", "started_at", "completed_at"],
      report.jobs.flatMap((job) =>
        job.steps.map((step) => [...key, job.name, step.number, step.name, step.status, step.started_at, step.completed_at]),
      ),
    ),
    ...inserts(
      db,
      "ci_substeps",
      ["run_id", "attempt", "job", "step", "action_id", "name", "depth", "completed_at", "duration_ms"],
      report.jobs.flatMap((job) =>
        uniqueIds(job.substeps ?? []).map((sub) => [
          ...key,
          job.name,
          sub.step,
          sub.action_id,
          sub.name,
          sub.depth,
          sub.completed_at,
          sub.duration_ms,
        ]),
      ),
    ),
    ...inserts(
      db,
      "ci_metrics",
      ["run_id", "attempt", "job", "source", "name", "value", "labels"],
      report.metrics.map((metric) => [
        ...key,
        metric.job,
        metric.source,
        metric.name,
        metric.value,
        sortedLabels(metric.labels),
      ]),
    ),
  ]);

  return new Response(null, { status: 204 });
}
