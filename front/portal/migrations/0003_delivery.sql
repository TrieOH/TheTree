-- What the GitHub webhook and the deploy hook need beyond 0001. TheTree runs
-- CI only on release tags, so pull requests carry no check state.

ALTER TABLE pull_requests DROP COLUMN checks;
ALTER TABLE pull_requests ADD COLUMN draft INTEGER NOT NULL DEFAULT 0;
ALTER TABLE pull_requests ADD COLUMN updated_at TEXT NOT NULL DEFAULT '';
CREATE INDEX pull_requests_state ON pull_requests (state, merged_at);

-- The Actions run that shipped it, for a link from the Production table.
ALTER TABLE deploys ADD COLUMN run_id INTEGER;

CREATE INDEX workflow_runs_by_workflow ON workflow_runs (workflow, queued_at);
