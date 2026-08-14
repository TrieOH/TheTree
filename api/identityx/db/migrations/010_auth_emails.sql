-- +goose Up
CREATE TABLE email_templates (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE,

    kind TEXT NOT NULL,
    CONSTRAINT chk_email_templates_kind CHECK (kind IN ('verify', 'reset')),

    subject TEXT NOT NULL,
    body TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_email_template_per_project UNIQUE (project_id, kind)
);

CREATE INDEX idx_email_templates_project_id ON email_templates(project_id);

CREATE TABLE action_tokens (
    jti UUID PRIMARY KEY,
    purpose TEXT NOT NULL,
    CONSTRAINT chk_action_tokens_purpose CHECK (purpose IN ('email_verify', 'password_reset')),

    actor_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,

    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_action_tokens_actor_id ON action_tokens(actor_id);
CREATE INDEX idx_action_tokens_expires_at ON action_tokens(expires_at);

-- +goose Down
DROP INDEX IF EXISTS idx_action_tokens_expires_at;
DROP INDEX IF EXISTS idx_action_tokens_actor_id;
DROP TABLE IF EXISTS action_tokens;
DROP INDEX IF EXISTS idx_email_templates_project_id;
DROP TABLE IF EXISTS email_templates;
