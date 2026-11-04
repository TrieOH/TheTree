-- Everything the portal shows, copied in from Linear, GitHub, the deploy
-- workflows, Cloudflare analytics, Grafana and Neon. Never read live.

CREATE TABLE issues (
  id TEXT PRIMARY KEY,
  identifier TEXT NOT NULL,
  title TEXT NOT NULL,
  state TEXT NOT NULL,
  priority INTEGER NOT NULL DEFAULT 0,
  estimate INTEGER,
  service TEXT,
  assignee TEXT,
  cycle_id TEXT,
  project_id TEXT,
  created_at TEXT NOT NULL,
  completed_at TEXT,
  updated_at TEXT NOT NULL
);
CREATE INDEX issues_cycle ON issues (cycle_id);
CREATE INDEX issues_service ON issues (service);

CREATE TABLE cycles (
  id TEXT PRIMARY KEY,
  number INTEGER NOT NULL,
  starts_at TEXT NOT NULL,
  ends_at TEXT NOT NULL
);

CREATE TABLE projects (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  progress REAL NOT NULL DEFAULT 0,
  target_date TEXT,
  updated_at TEXT NOT NULL
);

CREATE TABLE pull_requests (
  number INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  state TEXT NOT NULL,
  review_state TEXT,
  checks TEXT,
  opened_at TEXT NOT NULL,
  merged_at TEXT,
  merge_sha TEXT
);

CREATE TABLE workflow_runs (
  id INTEGER PRIMARY KEY,
  workflow TEXT NOT NULL,
  branch TEXT NOT NULL,
  sha TEXT NOT NULL,
  conclusion TEXT,
  attempt INTEGER NOT NULL DEFAULT 1,
  queued_at TEXT,
  started_at TEXT,
  completed_at TEXT,
  runner TEXT
);
CREATE INDEX workflow_runs_recent ON workflow_runs (workflow, branch, completed_at);

CREATE TABLE tags (
  name TEXT PRIMARY KEY,
  component TEXT NOT NULL,
  version TEXT NOT NULL,
  sha TEXT NOT NULL,
  message TEXT,
  created_at TEXT NOT NULL
);
CREATE INDEX tags_component ON tags (component, created_at);

CREATE TABLE deploys (
  id INTEGER PRIMARY KEY,
  component TEXT NOT NULL,
  tag TEXT NOT NULL,
  sha TEXT NOT NULL,
  actor TEXT NOT NULL,
  status TEXT NOT NULL,
  health_ok INTEGER,
  at TEXT NOT NULL
);
CREATE INDEX deploys_component ON deploys (component, at);

CREATE TABLE metric_rollups (
  service TEXT NOT NULL,
  metric TEXT NOT NULL,
  bucket TEXT NOT NULL,
  value REAL NOT NULL,
  PRIMARY KEY (service, metric, bucket)
);

CREATE TABLE usage_snapshots (
  resource TEXT NOT NULL,
  period TEXT NOT NULL,
  used REAL NOT NULL,
  included REAL NOT NULL,
  at TEXT NOT NULL,
  PRIMARY KEY (resource, period)
);

CREATE TABLE cold_starts (
  resource TEXT NOT NULL,
  kind TEXT NOT NULL,
  day TEXT NOT NULL,
  count INTEGER NOT NULL,
  p50_ms INTEGER,
  p95_ms INTEGER,
  PRIMARY KEY (resource, day)
);

-- Webhook redeliveries are dropped by (source, delivery id).
CREATE TABLE webhook_deliveries (
  source TEXT NOT NULL,
  delivery_id TEXT NOT NULL,
  received_at TEXT NOT NULL,
  PRIMARY KEY (source, delivery_id)
);
