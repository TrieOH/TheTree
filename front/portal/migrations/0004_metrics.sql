-- Phase 5: Cloudflare analytics and Neon, pulled every 30 minutes.
--
-- metric_rollups (0001) holds hourly traffic per service:
--   ('payssage', 'requests',   '2026-11-04T17:00:00Z', 12)
--   ('payssage', 'errors_5xx', '2026-11-04T17:00:00Z', 1)
-- usage_snapshots (0001) holds plan usage per billing cycle:
--   ('workers_requests', '2026-10-24', 52000, 10000000, at)

-- One row per Worker per billing cycle: the fronts' startup and the APIs'
-- edge Workers.
CREATE TABLE worker_stats (
  script TEXT NOT NULL,
  period TEXT NOT NULL,
  requests INTEGER NOT NULL,
  errors INTEGER NOT NULL,
  cpu_p50_us REAL,
  cpu_p99_us REAL,
  wall_p50_us REAL,
  wall_p99_us REAL,
  at TEXT NOT NULL,
  PRIMARY KEY (script, period)
);

-- One row per service: its Neon project, from the control-plane API only.
CREATE TABLE neon_projects (
  service TEXT PRIMARY KEY,
  project_id TEXT NOT NULL,
  period_start TEXT NOT NULL,
  period_end TEXT NOT NULL,
  storage_bytes INTEGER,
  active_seconds INTEGER,
  last_active_at TEXT,
  resumes INTEGER NOT NULL DEFAULT 0,
  resume_p50_ms INTEGER,
  resume_max_ms INTEGER,
  at TEXT NOT NULL
);
