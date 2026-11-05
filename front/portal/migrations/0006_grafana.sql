-- Phase 7: latency, slowest endpoints and error types from Grafana Cloud's
-- Prometheus query API (read-only token), pulled every 30 minutes with the
-- rest of the metrics. Measured inside each API's container, so a cold
-- start's wait is not in these numbers (cold_starts has it).
--
-- metric_rollups (0001) gains two metrics:
--   ('payssage', 'p95_ms',     '2026-11-04T17:00:00Z', 42.5)  hourly
--   ('payssage', 'p95_ms_24h', '2026-11-04T18:00:00Z', 38.1)  the 24h before that pull

-- One row per endpoint per UTC day. p95 is per route (the latency histogram
-- has no method label), requests per method and route.
CREATE TABLE endpoint_rollups (
  service TEXT NOT NULL,
  method TEXT NOT NULL,
  route TEXT NOT NULL,
  day TEXT NOT NULL,
  requests INTEGER NOT NULL,
  p95_ms REAL,
  PRIMARY KEY (service, method, route, day)
);

-- Problem answers per type slug per UTC day: slugs and counts only, never
-- messages (http_problems_total, lib/go/httpserver).
CREATE TABLE problem_rollups (
  service TEXT NOT NULL,
  problem TEXT NOT NULL,
  day TEXT NOT NULL,
  status INTEGER NOT NULL,
  count INTEGER NOT NULL,
  PRIMARY KEY (service, problem, day)
);
