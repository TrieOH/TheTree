-- +goose Up
CREATE TABLE namespaces (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uniq_namespace_name_per_user UNIQUE (owner_id, name)
);

CREATE INDEX idx_namespaces_owner_id ON namespaces (owner_id);

-- +goose Down
DROP INDEX IF EXISTS idx_namespaces_owner_id;
DROP TABLE IF EXISTS namespaces;