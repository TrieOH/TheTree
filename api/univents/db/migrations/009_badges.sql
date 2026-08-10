-- +goose Up
CREATE TABLE badge_templates (
    id             UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id     UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    ticket_type_id UUID REFERENCES ticket_types(id),
    name           VARCHAR(256) NOT NULL,
    design_data    JSONB NOT NULL DEFAULT '{}',
    origin         TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ
);

-- One template per (edition, ticket type, origin): at most one default, one
-- per ticket type, and one staff design per edition. NULLS NOT DISTINCT so
-- the (NULL ticket_type, NULL origin) default slot is unique.
CREATE UNIQUE INDEX idx_one_badge_template_per_scope
    ON badge_templates(edition_id, ticket_type_id, origin) NULLS NOT DISTINCT
    WHERE deleted_at IS NULL;

-- origin is either null (edition default or per-ticket-type scope) or 'staff'
-- (staff-only design, which never targets a ticket type).
ALTER TABLE badge_templates
    ADD CONSTRAINT chk_badge_template_scope_valid CHECK (
        origin IS NULL OR (origin = 'staff' AND ticket_type_id IS NULL)
    );

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
DROP TABLE IF EXISTS badge_templates;
