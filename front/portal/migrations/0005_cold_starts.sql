-- Phase 6: container cold starts, per service per day, from the
-- container_cold_starts Analytics Engine dataset the API Workers write.
-- cold_starts (0001) keys on (resource, day); resource is the service and
-- kind is 'container'. A failed start is a 5xx or no answer at all.
ALTER TABLE cold_starts ADD COLUMN failed INTEGER NOT NULL DEFAULT 0;
