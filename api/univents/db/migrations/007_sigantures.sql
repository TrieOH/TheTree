-- +goose Up
CREATE TABLE signatures (
    id                UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id        UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    created_by        UUID NOT NULL,
    signatory_name    VARCHAR(256) NOT NULL,
    signatory_title   VARCHAR(256),
    signatory_email   VARCHAR(256),
    signatory_user_id UUID,
    image_url         TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

CREATE TABLE signature_requests (
    id                UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id        UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    created_by        UUID NOT NULL,
    signatory_name    VARCHAR(256) NOT NULL,
    signatory_title   VARCHAR(256),
    signatory_email   VARCHAR(256),
    signatory_user_id UUID,
    idempotency_key   UUID NOT NULL,
    CONSTRAINT uniq_signature_requests_idempotency_key UNIQUE (idempotency_key),
    status            TEXT NOT NULL,
    status_reason     TEXT,
    expires_at        TIMESTAMPTZ NOT NULL,
    signature_id      UUID REFERENCES signatures(id),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT chk_signature_requests_status_valid CHECK (
        status IN ('pending', 'completed', 'expired', 'cancelled')
    )
);
-- +goose Down
DROP TABLE IF EXISTS signature_requests;
DROP TABLE IF EXISTS signatures;
