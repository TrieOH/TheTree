-- +goose Up
CREATE TABLE steps (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    form_id UUID NOT NULL REFERENCES forms(id)
        ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position_hint INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS steps;
