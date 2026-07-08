-- +goose Up
CREATE TABLE certification_templates (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id UUID NOT NULL REFERENCES editions(id)
    ON DELETE CASCADE,

    title TEXT NOT NULL,

    data JSONB NOT NULL,
    url TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE editions
    ADD COLUMN certification_template_id UUID REFERENCES certification_templates(id)
        ON DELETE SET NULL;

ALTER TABLE activities
    ADD COLUMN certification_template_id UUID REFERENCES certification_templates(id)
        ON DELETE SET NULL;

CREATE TABLE certifications (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,

    target_id UUID NOT NULL,
    target_type TEXT NOT NULL,
    CONSTRAINT chk_target_type_valid CHECK (
        target_type IN (
            'edition',
            'activity'
        )
    ),

    certified_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_certification_templates_edition_id ON certification_templates(edition_id);
CREATE INDEX idx_editions_certification_template_id ON editions(certification_template_id);
CREATE INDEX idx_activities_certification_template_id ON activities(certification_template_id);
CREATE INDEX idx_certifications_target ON certifications(target_type, target_id);

-- +goose Down
DROP INDEX IF EXISTS idx_certification_templates_edition_id;
DROP INDEX IF EXISTS idx_editions_certification_template_id;
DROP INDEX IF EXISTS idx_activities_certification_template_id;
DROP INDEX IF EXISTS idx_certifications_target;
DROP TABLE IF EXISTS certifications;
ALTER TABLE activities DROP COLUMN IF EXISTS certification_template_id;
ALTER TABLE editions DROP COLUMN IF EXISTS certification_template_id;
DROP TABLE IF EXISTS certification_templates;