-- +goose Up
CREATE TABLE answers (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    response_id UUID NOT NULL REFERENCES responses(id) ON DELETE CASCADE,
    field_id UUID REFERENCES fields(id) ON DELETE SET NULL,
    answer JSONB,
    answered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    CONSTRAINT uniq_answer_per_field_per_response UNIQUE (response_id, field_id)
);

-- +goose Down
DROP TABLE IF EXISTS answers;
