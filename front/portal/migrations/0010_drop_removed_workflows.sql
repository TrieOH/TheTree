-- warm-cache.yml (both of its names) was removed on 2026-11-05: the sticky
-- disks made it unnecessary. Its runs would otherwise stay on the CI tab.
DELETE FROM workflow_runs WHERE workflow IN ('Warm cache', 'Warm tool cache');
