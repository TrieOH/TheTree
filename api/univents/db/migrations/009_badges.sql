-- +goose Up
CREATE TABLE badge_templates (
    id             UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id     UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    ticket_type_id UUID REFERENCES ticket_types(id),
    name           VARCHAR(256) NOT NULL,
    design_data    JSONB NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_one_badge_template_per_ticket_type
    ON badge_templates(edition_id, ticket_type_id) NULLS NOT DISTINCT
    WHERE deleted_at IS NULL;
-- +goose Down
DROP TABLE IF EXISTS badge_templates;