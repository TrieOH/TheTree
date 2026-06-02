-- +goose Up
CREATE TABLE responders (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID,
    email TEXT NOT NULL,
    CONSTRAINT uniq_responder_email_on_system UNIQUE (email)
);

CREATE TABLE form_invites (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,

    -- named invite: link to existing responder
    responder_id UUID REFERENCES responders(id) ON DELETE SET NULL,
    -- named invite: email when responder record doesn't exist yet
    email TEXT,

    -- both NULL = open invite (anyone with the token)
    CONSTRAINT chk_invite_not_both_named
        CHECK (NOT (responder_id IS NOT NULL AND email IS NOT NULL)),

    expires_at TIMESTAMPTZ,
    accepted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE responses (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    invite_id UUID REFERENCES form_invites(id) ON DELETE SET NULL,
    responder_id UUID REFERENCES responders(id) ON DELETE SET NULL,
    email TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS responses;
DROP TABLE IF EXISTS form_invites;
DROP TABLE IF EXISTS responders;