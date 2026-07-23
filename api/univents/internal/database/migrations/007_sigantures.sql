-- +goose Up
CREATE TABLE signatures (
    id                UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id        UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    created_by        UUID NOT NULL,
    signatory_name    VARCHAR(256) NOT NULL,
    signatory_title   VARCHAR(256),
    signatory_email   VARCHAR(256),
    signatory_user_id UUID,
    image_url         TEXT,
    status            TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT chk_signatures_status_valid CHECK (
        status IN ('requested', 'ready', 'declined', 'expired')
        ),
    CONSTRAINT chk_signatures_ready_has_image CHECK (
        status != 'ready' OR image_url IS NOT NULL
    )
);
-- +goose Down
DROP TABLE IF EXISTS signatures;