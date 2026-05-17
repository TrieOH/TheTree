-- +goose Up
CREATE TABLE responders (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID,

    email TEXT NOT NULL,
    CONSTRAINT uniq_responder_email_on_system UNIQUE (email)
);

CREATE TABLE form_responders (
    form_id UUID NOT NULL REFERENCES forms(id)
        ON DELETE CASCADE,
    responder_id UUID NOT NULL REFERENCES responders(id)
        ON DELETE CASCADE,

    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (form_id, responder_id)
);

-- +goose Down
DROP TABLE IF EXISTS form_responders;
DROP TABLE IF EXISTS responders;