-- +goose Up
CREATE TABLE responses (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    form_id UUID NOT NULL REFERENCES forms(id)
        ON DELETE CASCADE,

    is_anonymous BOOLEAN DEFAULT FALSE,

    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS responses;
