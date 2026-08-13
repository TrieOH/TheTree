-- 006_webhook_events_status_detail.sql
-- +goose Up
-- Normalized failure/outcome detail (models.IntentStatusDetail, e.g.
-- "high_risk", "insufficient_funds") captured at parse time, so the
-- delivery worker can include it in the tenant envelope without
-- re-fetching the intent. NULL when the event carried no meaningful
-- detail (e.g. succeeded/pending).
ALTER TABLE webhook_events ADD COLUMN status_detail TEXT;

-- +goose Down
ALTER TABLE webhook_events DROP COLUMN status_detail;
