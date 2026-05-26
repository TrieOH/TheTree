-- +goose Up
CREATE TABLE answers (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    response_id UUID NOT NULL REFERENCES responses(id)
        ON DELETE CASCADE,
    field_id UUID NOT NULL REFERENCES fields(id)
        ON DELETE SET NULL,

    answer JSONB,

    answered_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS answers;
