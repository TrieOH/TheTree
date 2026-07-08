-- +goose Up
CREATE TABLE signatures (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id UUID NOT NULL REFERENCES editions(id)
        ON DELETE CASCADE,

    title TEXT NOT NULL,
    url TEXT NOT NULL,

    pos_x INT NOT NULL,
    pos_y INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS signatures;