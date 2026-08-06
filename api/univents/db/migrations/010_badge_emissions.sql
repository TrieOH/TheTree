-- +goose Up
CREATE TABLE badge_emissions (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id      UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL,
    origin          TEXT NOT NULL,
    registration_id UUID REFERENCES registrations(id),
    status          TEXT NOT NULL DEFAULT 'active',
    status_reason   TEXT,
    email_sent_at   TIMESTAMPTZ,
    emitted_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ,
    CONSTRAINT chk_badge_emissions_origin_valid CHECK (origin IN ('participant', 'staff')),
    CONSTRAINT chk_badge_emissions_status_valid CHECK (status IN ('active', 'revoked')),
    CONSTRAINT uniq_badge_emissions_edition_user_origin UNIQUE (edition_id, user_id, origin)
);

CREATE INDEX idx_badge_emissions_user ON badge_emissions(user_id);
CREATE INDEX idx_badge_emissions_edition ON badge_emissions(edition_id);
-- +goose Down
DROP TABLE IF EXISTS badge_emissions;
