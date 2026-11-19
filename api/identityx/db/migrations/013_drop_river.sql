-- +goose Up
-- River (the old job queue) is gone from the backends; these are the tables,
-- sequences, enum and function its own migrator created. Jobs now go through
-- Cloudflare Queues (lib/jobs). IF EXISTS: a database River never touched has
-- nothing to drop. lib/go/river stays; nothing here depends on it.
DROP TABLE IF EXISTS river_job, river_leader, river_migration, river_notification, river_queue, river_client_queue, river_client CASCADE;
DROP FUNCTION IF EXISTS river_job_state_in_bitmask(bit, river_job_state);
DROP FUNCTION IF EXISTS river_job_notify();
DROP TYPE IF EXISTS river_job_state;
DROP SEQUENCE IF EXISTS river_job_id_seq, river_notification_id_seq;

-- +goose Down
-- Not recreated: River rebuilds its own schema with its migrator if it is ever
-- brought back.
SELECT 1;
