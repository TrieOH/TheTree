-- +goose Up
ALTER TABLE badge_templates ADD COLUMN origin TEXT;

-- The old (edition_id, ticket_type_id) unique index treated every NULL
-- ticket_type_id row as the same slot, so a staff template (ticket_type_id
-- NULL, origin 'staff') could not coexist with the edition default (both
-- ticket_type_id NULL). The new index scopes by origin too: one template per
-- (edition, ticket_type, origin) — i.e. at most one default, one per ticket
-- type, and one staff design per edition.
DROP INDEX idx_one_badge_template_per_ticket_type;

CREATE UNIQUE INDEX idx_one_badge_template_per_scope
    ON badge_templates(edition_id, ticket_type_id, origin) NULLS NOT DISTINCT
    WHERE deleted_at IS NULL;

-- origin is either null (edition default or per-ticket-type scope) or 'staff'
-- (staff-only design, which never targets a ticket type).
ALTER TABLE badge_templates
    ADD CONSTRAINT chk_badge_template_scope_valid CHECK (
        origin IS NULL OR (origin = 'staff' AND ticket_type_id IS NULL)
    );
-- +goose Down
ALTER TABLE badge_templates DROP CONSTRAINT chk_badge_template_scope_valid;

DROP INDEX idx_one_badge_template_per_scope;

CREATE UNIQUE INDEX idx_one_badge_template_per_ticket_type
    ON badge_templates(edition_id, ticket_type_id) NULLS NOT DISTINCT
    WHERE deleted_at IS NULL;

ALTER TABLE badge_templates DROP COLUMN origin;
