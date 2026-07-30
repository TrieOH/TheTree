-- +goose Up
CREATE TABLE certification_templates (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id   UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    kind         TEXT NOT NULL,
    CONSTRAINT chk_certification_templates_kind_valid CHECK (
        kind IN ('edition_attendance', 'program_attendance')
    ),
    name         VARCHAR(256) NOT NULL,
    description  TEXT,
    design_data  JSONB NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ
);

CREATE TABLE certification_template_programs (
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    template_id UUID NOT NULL REFERENCES certification_templates(id) ON DELETE CASCADE,
    program_id  UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_cert_template_program_unique
    ON certification_template_programs(template_id, program_id);

CREATE TABLE certifications (
    id                UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id        UUID NOT NULL REFERENCES editions(id),
    template_id       UUID REFERENCES certification_templates(id),
    registration_id   UUID NOT NULL REFERENCES registrations(id),
    user_id           UUID NOT NULL,
    program_id        UUID REFERENCES programs(id),
    verification_hash TEXT NOT NULL,
    valid             BOOLEAN NOT NULL DEFAULT true,
    invalid_reason    TEXT,
    email_sent        BOOLEAN NOT NULL DEFAULT false,
    issued_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_one_cert_per_registration_per_template
    ON certifications(template_id, registration_id);

CREATE UNIQUE INDEX idx_one_edition_cert_per_registration
    ON certifications(registration_id)
    WHERE template_id IS NULL AND program_id IS NULL;

CREATE UNIQUE INDEX idx_one_cert_per_user_per_program
    ON certifications(user_id, program_id)
    WHERE program_id IS NOT NULL;

CREATE UNIQUE INDEX idx_certifications_verification_hash
    ON certifications(verification_hash);

CREATE TABLE cert_emission_errors (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id   UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL,
    template_id  UUID REFERENCES certification_templates(id),
    program_id   UUID REFERENCES programs(id),
    error_message TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE IF EXISTS cert_emission_errors;
DROP TABLE IF EXISTS certifications;
DROP TABLE IF EXISTS certification_template_programs;
DROP TABLE IF EXISTS certification_templates;
