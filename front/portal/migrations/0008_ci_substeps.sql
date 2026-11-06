-- What happens inside a step that runs a composite action (Setup Go is
-- setup-go + two sticky disks + ...). The Actions API only times top-level
-- steps; the report job reads each job's log, where the runner marks every
-- action it runs (`##[end-action id=...;duration_ms=N]`), and sends those too.

CREATE TABLE ci_substeps (
  run_id INTEGER NOT NULL,
  attempt INTEGER NOT NULL,
  job TEXT NOT NULL,
  -- the top-level step it ran inside (ci_steps.number), by time
  step INTEGER,
  -- the runner's id for it, e.g. __self_2.__actions_setup-node; dots nest
  action_id TEXT NOT NULL,
  -- what the log shows for it, e.g. "Toolchain" or "Run pnpm install ..."
  name TEXT NOT NULL,
  depth INTEGER NOT NULL,
  completed_at TEXT NOT NULL,
  duration_ms INTEGER NOT NULL,
  PRIMARY KEY (run_id, attempt, job, action_id)
);
