-- What the Linear webhook and the hourly catch-up sync need beyond 0001.
-- The service an issue belongs to is its Linear project, joined at read
-- time, so a renamed project needs no backfill; `service` goes.

DROP INDEX issues_service;
ALTER TABLE issues DROP COLUMN service;
ALTER TABLE issues ADD COLUMN state_type TEXT NOT NULL DEFAULT 'unstarted';
ALTER TABLE issues ADD COLUMN labels TEXT NOT NULL DEFAULT '[]';
ALTER TABLE issues ADD COLUMN url TEXT;
CREATE INDEX issues_project ON issues (project_id);
CREATE INDEX issues_state_type ON issues (state_type, completed_at);

-- Project status type (backlog, planned, started, paused, completed,
-- canceled). Progress is counted from the issues, so it stays unused.
ALTER TABLE projects ADD COLUMN state TEXT;

ALTER TABLE cycles ADD COLUMN completed_at TEXT;

-- Where each pull left off, e.g. ('linear', '2026-11-04T16:00:00.000Z').
CREATE TABLE sync_state (
  source TEXT PRIMARY KEY,
  synced_at TEXT NOT NULL
);
