-- 007_webhook_events_dedupe_by_event_type.sql
-- +goose Up
-- MP fires multiple webhook notifications per payment (payment.created,
-- payment.pending, payment.updated, ...), each representing a different
-- status. The old dedupe on (provider, external_id) let only the FIRST
-- event through — a later, more important event (payment.succeeded /
-- payment.rejected) hit the unique constraint and was silently dropped in
-- receive.go's conflict branch, so the final status never reached tenant
-- endpoints (purchases expired with no confirmation). Dedupe per
-- (provider, external_id, event_type) instead: identical redeliveries
-- still conflict (the first was already dispatched), distinct statuses
-- each record and dispatch.
DROP INDEX IF EXISTS uniq_webhook_events_external_id;
CREATE UNIQUE INDEX uniq_webhook_events_external_id_type
    ON webhook_events (provider, external_id, event_type)
    WHERE external_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS uniq_webhook_events_external_id_type;
CREATE UNIQUE INDEX uniq_webhook_events_external_id
    ON webhook_events (provider, external_id)
    WHERE external_id IS NOT NULL;
