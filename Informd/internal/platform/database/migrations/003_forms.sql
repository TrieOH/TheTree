-- +goose Up
CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE,
    owner_id UUID NOT NULL,
    title TEXT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    --current_version_id UUID REFERENCES versions(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    opened_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,

    CONSTRAINT chk_forms_status CHECK (status IN ('draft', 'open', 'closed', 'archived')),
    CONSTRAINT chk_forms_valid_status_state CHECK (
        (status = 'open' AND opened_at IS NOT NULL)
        OR (status = 'closed' AND closed_at IS NOT NULL)
        OR (status = 'archived' AND archived_at IS NOT NULL)
        OR (status = 'draft')
    )
);

CREATE UNIQUE INDEX uniq_idx_forms_title_project
    ON forms (title, project_id);

--CREATE INDEX idx_forms_current_version_id
--    ON forms (current_version_id);

CREATE INDEX idx_forms_owner_id ON forms (owner_id);

CREATE TABLE versions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    form_id UUID NOT NULL REFERENCES forms(id)
         ON DELETE CASCADE,

    version INT NOT NULL,
    CONSTRAINT chk_version_gt_zero CHECK (version > 0),

    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    CONSTRAINT chk_versions_status CHECK (status IN ('draft', 'active', 'outdated')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    activated_at TIMESTAMPTZ,
    outdated_at TIMESTAMPTZ
);

CREATE INDEX idx_versions_forms_id ON versions(form_id);

CREATE UNIQUE INDEX uniq_idx_version_number
    ON versions (form_id, version);

CREATE UNIQUE INDEX one_version_draft_per_form
    ON versions (form_id)
    WHERE status = 'draft';

CREATE UNIQUE INDEX one_version_active_per_form
    ON versions (form_id)
    WHERE status = 'active';

CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Stable logical identity of a field within a form
    stable_id UUID NOT NULL DEFAULT uuidv7(),
    version_id UUID NOT NULL REFERENCES versions(id)
        ON DELETE CASCADE,

    key VARCHAR(64) NOT NULL,
    CONSTRAINT uniq_one_key_per_version UNIQUE (version_id, key),

    type VARCHAR(32) NOT NULL DEFAULT 'string',
    CONSTRAINT chk_fields_type CHECK (type IN (
        'string',
        'email',
        'int',
        'float',
        'bool',
        'select'
    )),

    -- Access control after initial completion
    owner VARCHAR(32) NOT NULL DEFAULT 'user',
    CONSTRAINT chk_fields_owner CHECK (owner IN ('user', 'admin')),
    -- Means it can be changed after its set
    mutable BOOLEAN NOT NULL DEFAULT true,

    title VARCHAR(255) NOT NULL,
    description TEXT,
    placeholder JSONB,
    select_behaviour VARCHAR(32),
    select_options JSONB,
    CONSTRAINT chk_fields_select_behaviour CHECK (type <> 'select' OR select_behaviour IN ('checkbox', 'radio')),
    CONSTRAINT chk_fields_select_options CHECK (type <> 'select' OR (select_options IS NOT NULL AND jsonb_typeof(select_options) = 'array')),

    default_value JSONB,

    required BOOLEAN NOT NULL DEFAULT false,

    position INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- A field identity appears at most once per version
    CONSTRAINT uniq_one_stable_per_version UNIQUE (version_id, stable_id),

    CONSTRAINT chk_fields_key_format CHECK (key ~ '^[a-z_][a-z0-9_]*$')
);

CREATE INDEX idx_fields_version_id
    ON fields(version_id);

-- +goose Down
DROP INDEX IF EXISTS idx_fields_version_id;
DROP TABLE IF EXISTS fields;
DROP INDEX IF EXISTS one_version_draft_per_form;
DROP INDEX IF EXISTS uniq_idx_version_number;
DROP INDEX IF EXISTS idx_versions_forms_id;
DROP TABLE IF EXISTS versions;
DROP INDEX IF EXISTS idx_forms_owner_id;
--DROP INDEX IF EXISTS idx_forms_current_version_id;
DROP INDEX IF EXISTS uniq_idx_forms_title_project;
DROP TABLE IF EXISTS forms;
