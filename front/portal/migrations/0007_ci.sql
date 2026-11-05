-- What every release run spent its time on, sent by deploy.yml's report job
-- to POST /webhooks/ci: the run, each job (with how long it waited for a
-- runner), each step's start and end, and the numbers the actions recorded
-- with .github/scripts/metric (cache hits, test counts, findings, sizes).
-- A re-run is its own attempt.

CREATE TABLE ci_runs (
  run_id INTEGER NOT NULL,
  attempt INTEGER NOT NULL,
  component TEXT NOT NULL,
  tag TEXT NOT NULL,
  sha TEXT NOT NULL,
  type TEXT NOT NULL,
  lane TEXT NOT NULL,
  status TEXT NOT NULL,
  stage TEXT NOT NULL,
  started_at TEXT,
  completed_at TEXT,
  PRIMARY KEY (run_id, attempt)
);
CREATE INDEX ci_runs_by_type ON ci_runs (type, started_at);

CREATE TABLE ci_jobs (
  run_id INTEGER NOT NULL,
  attempt INTEGER NOT NULL,
  name TEXT NOT NULL,
  status TEXT,
  queued_at TEXT,
  started_at TEXT,
  completed_at TEXT,
  runner TEXT,
  PRIMARY KEY (run_id, attempt, name)
);

CREATE TABLE ci_steps (
  run_id INTEGER NOT NULL,
  attempt INTEGER NOT NULL,
  job TEXT NOT NULL,
  number INTEGER NOT NULL,
  name TEXT NOT NULL,
  status TEXT,
  started_at TEXT,
  completed_at TEXT,
  PRIMARY KEY (run_id, attempt, job, number)
);

CREATE TABLE ci_metrics (
  id INTEGER PRIMARY KEY,
  run_id INTEGER NOT NULL,
  attempt INTEGER NOT NULL,
  job TEXT NOT NULL,
  source TEXT NOT NULL,
  name TEXT NOT NULL,
  value REAL NOT NULL,
  -- JSON object, keys sorted, '{}' without labels
  labels TEXT NOT NULL DEFAULT '{}'
);
CREATE INDEX ci_metrics_by_run ON ci_metrics (run_id, attempt);

-- Where a release stopped: parse, checks (blocked before shipping anything)
-- or publish. NULL for reports from before the column.
ALTER TABLE deploys ADD COLUMN stage TEXT;
