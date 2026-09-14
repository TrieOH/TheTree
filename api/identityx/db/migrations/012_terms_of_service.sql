-- +goose Up
-- Terms of service: one current document per project (versioned in place),
-- plus the acceptance ledger proving which actor agreed to which version.
--
-- A project without a row in terms_of_service has no terms (v0): registration
-- asks for no acceptance and no notification is ever sent. Introducing the
-- first row (v1) is itself a version event — existing users are notified
-- exactly like on any later update, they simply migrate from "no terms"
-- to v1.
CREATE TABLE terms_of_service (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE,

    version INT NOT NULL DEFAULT 1,
    CONSTRAINT chk_terms_of_service_version CHECK (version >= 1),

    content TEXT NOT NULL,
    CONSTRAINT chk_terms_of_service_content CHECK (length(content) > 0),

    -- When the version takes effect for existing users. Defaults to the
    -- 30-day notice period customary under the Brazilian consumer legal
    -- framework (CDC) and expected by LGPD art. 9 §6 notification duty.
    effective_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '30 days'),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_terms_of_service_per_project UNIQUE (project_id)
);

CREATE INDEX idx_terms_of_service_project_id ON terms_of_service(project_id);

-- Consent ledger. The controller bears the burden of proving consent
-- (LGPD art. 8 §2), so each row pins the accepted version, the SHA-256
-- hash of the content at acceptance time, and the timestamp. Source tells
-- HOW the consent was given: 'clickwrap' is an explicit box-check
-- (registration or the accept endpoint), 'continued_use' is the first
-- login after the version's effective date — weaker but documented
-- evidence. Accepting the same version twice is idempotent (the unique
-- constraint collapses it; the repo refreshes accepted_at on conflict).
CREATE TABLE tos_acceptances (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    actor_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE,

    tos_version INT NOT NULL,
    content_hash TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'clickwrap',
    accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_tos_acceptance_per_actor_version UNIQUE (actor_id, project_id, tos_version),
    CONSTRAINT chk_tos_acceptances_source CHECK (source IN ('clickwrap', 'continued_use'))
);

CREATE INDEX idx_tos_acceptances_project_id ON tos_acceptances(project_id);
CREATE INDEX idx_tos_acceptances_actor_id ON tos_acceptances(actor_id);

-- The ToS-change notification is a third email template kind: projects may
-- override its copy (e.g. to write it in Portuguese) like verify/reset.
ALTER TABLE email_templates
    DROP CONSTRAINT chk_email_templates_kind;
ALTER TABLE email_templates
    ADD CONSTRAINT chk_email_templates_kind CHECK (kind IN ('verify', 'reset', 'tos'));

-- +goose Down
ALTER TABLE email_templates
    DROP CONSTRAINT chk_email_templates_kind;
ALTER TABLE email_templates
    ADD CONSTRAINT chk_email_templates_kind CHECK (kind IN ('verify', 'reset'));
DELETE FROM email_templates WHERE kind = 'tos';

DROP INDEX IF EXISTS idx_tos_acceptances_actor_id;
DROP INDEX IF EXISTS idx_tos_acceptances_project_id;
DROP TABLE IF EXISTS tos_acceptances;
DROP INDEX IF EXISTS idx_terms_of_service_project_id;
DROP TABLE IF EXISTS terms_of_service;
