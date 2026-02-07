-- +goose Up

CREATE TABLE token_reuse_list (
    jit UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS token_reuse_list;
