-- +goose Up
-- First-class caller correlation fields on intents (issue: univents refund plan A7).
-- external_id    = the caller's order id (univents purchase id)
-- external_group = the caller's grouping id (univents edition id)
-- Existing intents get backfilled by hand in prod (single edition); see the
-- refund plan's A7 — deliberately not a data migration here.
ALTER TABLE intents ADD COLUMN external_id TEXT;
ALTER TABLE intents ADD COLUMN external_group TEXT;

CREATE INDEX idx_intents_external_group ON intents (external_group);

-- +goose Down
DROP INDEX IF EXISTS idx_intents_external_group;
ALTER TABLE intents DROP COLUMN external_group;
ALTER TABLE intents DROP COLUMN external_id;
