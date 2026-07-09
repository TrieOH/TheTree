-- +goose Up
-- Add verification hash column to certifications.
-- NULL for pre-existing certs (they predate verification), always set for new ones.
ALTER TABLE certifications
    ADD COLUMN hash TEXT;

ALTER TABLE certifications
    ADD CONSTRAINT uniq_hash_per_certification UNIQUE(hash);

CREATE INDEX idx_certifications_hash ON certifications (hash)
    WHERE hash IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_certifications_hash;
ALTER TABLE certifications DROP CONSTRAINT IF EXISTS uniq_hash_per_certification;
ALTER TABLE certifications DROP COLUMN IF EXISTS hash;
