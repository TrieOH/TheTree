// GET /api/ci/*: the raw release-run data in D1, public and filterable, so a
// question about CI is one request instead of digging through GitHub logs.
// Every list takes the filters below, newest first, `limit` (default 200, at
// most 5000) and `format=csv`. GET /api/ci describes it all.

const MAX_LIMIT = 5000;
const DEFAULT_LIMIT = 200;

const ms = (from: string, to: string) =>
  `CAST(ROUND((julianday(${to}) - julianday(${from})) * 86400000) AS INTEGER)`;

/** Filters every list shares: they narrow the runs. */
const RUN_FILTERS = {
  type: "service, front, sdk or pages",
  component: "the release artifact, e.g. identityx or univents-ui",
  tag: "an exact tag, e.g. identityx/v0.48.4",
  status: "success, failure or cancelled",
  stage: "parse, checks or publish (where it stopped)",
  lane: "full or express",
  since: "ISO time; runs started at or after it",
  until: "ISO time; runs started before it",
} as const;

const ENDPOINTS = {
  "/api/ci/runs": {
    description: "One row per release run (attempt).",
    filters: RUN_FILTERS,
  },
  "/api/ci/runs/{run_id}": {
    description: "One run in full: its jobs with steps and substeps, and its metrics. ?attempt=N, default the latest.",
    filters: {},
  },
  "/api/ci/steps": {
    description: "One row per top-level step per run, with duration_ms. Waiting for a runner is job-level: see queue_ms in /api/ci/runs/{run_id}.",
    filters: { ...RUN_FILTERS, job: "checks or publish (or parse, report)", name: "the step's name, e.g. Setup Go", status_step: "the step's own result, e.g. skipped" },
  },
  "/api/ci/substeps": {
    description: "One row per action inside a composite step (from the job logs), with the step it ran in.",
    filters: { ...RUN_FILTERS, job: "checks or publish", step: "the top-level step's name, e.g. Setup Node + pnpm", name: "the action's name, e.g. Toolchain" },
  },
  "/api/ci/metrics": {
    description: "One row per number an action recorded (.github/scripts/metric).",
    filters: { ...RUN_FILTERS, job: "checks or publish", source: "e.g. setup/node", name: "e.g. pnpm_downloaded", "label.<key>": "a label's value, e.g. label.mode=short" },
  },
} as const;

type Row = Record<string, unknown>;

function respond(rows: unknown, format: string | null, status = 200): Response {
  const headers = {
    "Cache-Control": "no-store",
    "Access-Control-Allow-Origin": "*",
  };
  if (format === "csv" && Array.isArray(rows)) {
    return new Response(toCsv(rows as Row[]), { status, headers: { ...headers, "Content-Type": "text/csv; charset=utf-8" } });
  }
  return Response.json(rows, { status, headers });
}

function toCsv(rows: Row[]): string {
  if (rows.length === 0) return "";
  const columns = Object.keys(rows[0]);
  const cell = (value: unknown) => {
    const text = value === null || value === undefined ? "" : typeof value === "object" ? JSON.stringify(value) : String(value);
    return /[",\n]/.test(text) ? `"${text.replaceAll('"', '""')}"` : text;
  };
  return [columns.join(","), ...rows.map((row) => columns.map((c) => cell(row[c])).join(","))].join("\n") + "\n";
}

/** WHERE clauses over the run (alias r) from the query string. */
function runWhere(params: URLSearchParams): { sql: string[]; binds: unknown[] } {
  const sql: string[] = [];
  const binds: unknown[] = [];
  for (const column of ["type", "component", "tag", "status", "stage", "lane"] as const) {
    const value = params.get(column);
    if (value !== null) {
      sql.push(`r.${column} = ?`);
      binds.push(value);
    }
  }
  if (params.get("since")) {
    sql.push("r.started_at >= ?");
    binds.push(params.get("since"));
  }
  if (params.get("until")) {
    sql.push("r.started_at < ?");
    binds.push(params.get("until"));
  }
  return { sql, binds };
}

/** The job filter matches the job's last segment: "service / checks" → checks. */
function jobWhere(params: URLSearchParams, column: string, where: { sql: string[]; binds: unknown[] }) {
  const job = params.get("job");
  if (job !== null) {
    where.sql.push(`(${column} = ? OR ${column} LIKE ?)`);
    where.binds.push(job, `% / ${job}`);
  }
}

function limitOf(params: URLSearchParams): number {
  const limit = Number(params.get("limit") ?? DEFAULT_LIMIT);
  return Number.isInteger(limit) && limit > 0 ? Math.min(limit, MAX_LIMIT) : DEFAULT_LIMIT;
}

const clause = (sql: string[]) => (sql.length ? `WHERE ${sql.join(" AND ")}` : "");

async function list(db: D1Database, query: string, binds: unknown[]): Promise<Row[]> {
  return (await db.prepare(query).bind(...binds).all<Row>()).results;
}

async function runs(db: D1Database, params: URLSearchParams): Promise<Row[]> {
  const where = runWhere(params);
  return list(
    db,
    `SELECT r.run_id, r.attempt, r.component, r.tag, r.sha, r.type, r.lane, r.status, r.stage,
            r.started_at, r.completed_at, ${ms("r.started_at", "r.completed_at")} AS duration_ms
     FROM ci_runs r ${clause(where.sql)} ORDER BY r.started_at DESC LIMIT ?`,
    [...where.binds, limitOf(params)],
  );
}

async function steps(db: D1Database, params: URLSearchParams): Promise<Row[]> {
  const where = runWhere(params);
  jobWhere(params, "s.job", where);
  if (params.get("name") !== null) {
    where.sql.push("s.name = ?");
    where.binds.push(params.get("name"));
  }
  if (params.get("status_step") !== null) {
    where.sql.push("s.status = ?");
    where.binds.push(params.get("status_step"));
  }
  return list(
    db,
    `SELECT r.run_id, r.attempt, r.component, r.tag, r.type, r.lane, r.status AS run_status,
            s.job, s.number, s.name, s.status, s.started_at, s.completed_at,
            ${ms("s.started_at", "s.completed_at")} AS duration_ms
     FROM ci_steps s JOIN ci_runs r USING (run_id, attempt)
     ${clause(where.sql)} ORDER BY r.started_at DESC, s.job, s.number LIMIT ?`,
    [...where.binds, limitOf(params)],
  );
}

async function substeps(db: D1Database, params: URLSearchParams): Promise<Row[]> {
  const where = runWhere(params);
  jobWhere(params, "u.job", where);
  if (params.get("step") !== null) {
    where.sql.push("s.name = ?");
    where.binds.push(params.get("step"));
  }
  if (params.get("name") !== null) {
    where.sql.push("u.name = ?");
    where.binds.push(params.get("name"));
  }
  return list(
    db,
    `SELECT r.run_id, r.attempt, r.component, r.tag, r.type, r.lane,
            u.job, u.step AS step_number, s.name AS step, u.action_id, u.name, u.depth,
            u.completed_at, u.duration_ms
     FROM ci_substeps u JOIN ci_runs r USING (run_id, attempt)
     LEFT JOIN ci_steps s ON s.run_id = u.run_id AND s.attempt = u.attempt AND s.job = u.job AND s.number = u.step
     ${clause(where.sql)} ORDER BY r.started_at DESC, u.job, u.completed_at LIMIT ?`,
    [...where.binds, limitOf(params)],
  );
}

async function metrics(db: D1Database, params: URLSearchParams): Promise<Row[]> {
  const where = runWhere(params);
  jobWhere(params, "m.job", where);
  for (const column of ["source", "name"] as const) {
    if (params.get(column) !== null) {
      where.sql.push(`m.${column} = ?`);
      where.binds.push(params.get(column));
    }
  }
  for (const [key, value] of params) {
    if (!key.startsWith("label.")) continue;
    where.sql.push("json_extract(m.labels, ?) = ?");
    where.binds.push(`$."${key.slice("label.".length).replaceAll('"', "")}"`, value);
  }
  const rows = await list(
    db,
    `SELECT r.run_id, r.attempt, r.component, r.tag, r.type, r.lane,
            m.job, m.source, m.name, m.value, m.labels
     FROM ci_metrics m JOIN ci_runs r USING (run_id, attempt)
     ${clause(where.sql)} ORDER BY r.started_at DESC, m.id LIMIT ?`,
    [...where.binds, limitOf(params)],
  );
  return rows.map((row) => ({ ...row, labels: JSON.parse(row.labels as string) as Record<string, string> }));
}

async function oneRun(db: D1Database, runId: number, attemptParam: string | null): Promise<Row | null> {
  const run = await db
    .prepare(
      `SELECT run_id, attempt, component, tag, sha, type, lane, status, stage, started_at, completed_at,
              ${ms("started_at", "completed_at")} AS duration_ms
       FROM ci_runs WHERE run_id = ? ${attemptParam ? "AND attempt = ?" : ""} ORDER BY attempt DESC LIMIT 1`,
    )
    .bind(...(attemptParam ? [runId, Number(attemptParam)] : [runId]))
    .first<Row>();
  if (!run) return null;
  const key = [run.run_id, run.attempt];
  const [jobs, stepRows, subRows, metricRows] = await db.batch<Row>([
    db
      .prepare(
        `SELECT name, status, runner, queued_at, started_at, completed_at,
                ${ms("queued_at", "started_at")} AS queue_ms, ${ms("started_at", "completed_at")} AS duration_ms
         FROM ci_jobs WHERE run_id = ? AND attempt = ? ORDER BY started_at`,
      )
      .bind(...key),
    db
      .prepare(
        `SELECT job, number, name, status, started_at, completed_at, ${ms("started_at", "completed_at")} AS duration_ms
         FROM ci_steps WHERE run_id = ? AND attempt = ? ORDER BY job, number`,
      )
      .bind(...key),
    db
      .prepare(
        `SELECT job, step, action_id, name, depth, completed_at, duration_ms
         FROM ci_substeps WHERE run_id = ? AND attempt = ? ORDER BY job, completed_at`,
      )
      .bind(...key),
    db.prepare("SELECT job, source, name, value, labels FROM ci_metrics WHERE run_id = ? AND attempt = ? ORDER BY id").bind(...key),
  ]);
  const subs = subRows.results;
  return {
    ...run,
    jobs: jobs.results.map((job) => ({
      ...job,
      steps: stepRows.results
        .filter((step) => step.job === job.name)
        .map(({ job: _job, ...step }) => ({
          ...step,
          substeps: subs.filter((sub) => sub.job === job.name && sub.step === step.number).map(({ job: _j, step: _s, ...sub }) => sub),
        })),
    })),
    metrics: metricRows.results.map((row) => ({ ...row, labels: JSON.parse(row.labels as string) as Record<string, string> })),
  };
}

export async function handleCiApi(request: Request, db: D1Database, path: string): Promise<Response> {
  const params = new URL(request.url).searchParams;
  const format = params.get("format");

  if (path === "/api/ci") {
    return respond({
      about: "TheTree release runs (deploy.yml), step by step. Public and read-only.",
      common: { limit: `rows, default ${DEFAULT_LIMIT}, at most ${MAX_LIMIT}`, format: "json (default) or csv" },
      endpoints: ENDPOINTS,
    }, null);
  }
  if (path === "/api/ci/runs") return respond(await runs(db, params), format);
  if (path === "/api/ci/steps") return respond(await steps(db, params), format);
  if (path === "/api/ci/substeps") return respond(await substeps(db, params), format);
  if (path === "/api/ci/metrics") return respond(await metrics(db, params), format);

  const match = /^\/api\/ci\/runs\/(\d+)$/.exec(path);
  if (match) {
    const run = await oneRun(db, Number(match[1]), params.get("attempt"));
    return run ? respond(run, null) : respond({ error: "no such run" }, null, 404);
  }
  return respond({ error: "not found", see: "/api/ci" }, null, 404);
}
