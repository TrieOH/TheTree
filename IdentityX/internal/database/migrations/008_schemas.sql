-- +goose Up

-- =========================
-- ENUM TYPES
-- =========================

CREATE TYPE schema_type AS ENUM ('core', 'context', 'sub-context');
CREATE TYPE schema_status AS ENUM ('draft', 'published', 'archived');

CREATE TYPE schema_version_status AS ENUM ('draft', 'published', 'archived');

CREATE TYPE field_type AS ENUM (
    'string',
    'int',
    'select',
    'radio',
    'checkbox',
    'bool'
);

CREATE TYPE field_owner AS ENUM (
    'system',
    'admin',
    'user'
);

CREATE TYPE rule_operator AS ENUM (
    'equals',
    'not_equals',
    'in',
    'not_in',
    'exists',
    'not_exists'
);

-- =========================
-- SCHEMAS
-- =========================

CREATE TABLE schemas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL DEFAULT 'Unnamed Schema',
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    flow_id VARCHAR(63) NOT NULL,
    type schema_type NOT NULL,
    current_version_id UUID,
    status schema_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (project_id, type, flow_id)
);

-- =========================
-- SCHEMA VERSIONS
-- =========================

CREATE TABLE schema_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schema_id UUID NOT NULL REFERENCES schemas(id) ON DELETE CASCADE,
    version INT NOT NULL,
    status schema_version_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (schema_id, version)
);

CREATE INDEX idx_schema_versions_schema_id
    ON schema_versions(schema_id);

CREATE UNIQUE INDEX uniq_published_schema_versions
    ON schema_versions (schema_id, version)
    WHERE status = 'published';

-- =========================
-- FIELDS
-- =========================

CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schema_version_id UUID NOT NULL REFERENCES schema_versions(id) ON DELETE CASCADE,

    key VARCHAR(63) NOT NULL,                -- machine identifier
    type field_type NOT NULL DEFAULT 'string',
    owner field_owner NOT NULL DEFAULT 'system',

    title VARCHAR(255) NOT NULL,
    description TEXT,
    placeholder VARCHAR(255),

    required BOOLEAN NOT NULL DEFAULT false, -- base requirement
    mutable BOOLEAN NOT NULL DEFAULT true,

    default_value JSONB,
    position INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (schema_version_id, key)
);

CREATE INDEX idx_fields_schema_version_id
    ON fields(schema_version_id);

-- =========================
-- FIELD OPTIONS
-- =========================

CREATE TABLE field_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    field_id UUID NOT NULL REFERENCES fields(id) ON DELETE CASCADE,

    value VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    position INT NOT NULL,

    UNIQUE (field_id, value)
);

CREATE INDEX idx_field_options_field_id
    ON field_options(field_id);

-- =========================
-- FIELD VISIBILITY RULES
-- =========================

CREATE TABLE field_visibility_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    field_id UUID NOT NULL
        REFERENCES fields(id)
            ON DELETE CASCADE,

    depends_on_field_id UUID NOT NULL
        REFERENCES fields(id)
            ON DELETE CASCADE,

    operator rule_operator NOT NULL,
    value JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (field_id <> depends_on_field_id)
);

CREATE INDEX idx_visibility_rules_field_id
    ON field_visibility_rules(field_id);

CREATE INDEX idx_visibility_rules_depends_on
    ON field_visibility_rules(depends_on_field_id);

-- =========================
-- FIELD REQUIRED RULES
-- =========================

CREATE TABLE field_required_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    field_id UUID NOT NULL
        REFERENCES fields(id)
            ON DELETE CASCADE,

    depends_on_field_id UUID NOT NULL
        REFERENCES fields(id)
            ON DELETE CASCADE,

    operator rule_operator NOT NULL,
    value JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (field_id <> depends_on_field_id)
);

CREATE INDEX idx_required_rules_field_id
    ON field_required_rules(field_id);

CREATE INDEX idx_required_rules_depends_on
    ON field_required_rules(depends_on_field_id);

-- =========================
-- CURRENT VERSION FK
-- =========================

ALTER TABLE schemas
    ADD CONSTRAINT fk_current_schema_version
        FOREIGN KEY (current_version_id)
            REFERENCES schema_versions(id)
            ON DELETE SET NULL;

-- +goose Down
ALTER TABLE schemas
DROP CONSTRAINT IF EXISTS fk_current_schema_version;

DROP TABLE IF EXISTS field_required_rules;
DROP TABLE IF EXISTS field_visibility_rules;
DROP TYPE IF EXISTS rule_operator;

DROP TABLE IF EXISTS field_options;

DROP TABLE IF EXISTS fields;
DROP TYPE IF EXISTS field_owner;
DROP TYPE IF EXISTS field_type;

DROP TABLE IF EXISTS schema_versions;
DROP TYPE IF EXISTS schema_version_status;

DROP TABLE IF EXISTS schemas;
DROP TYPE IF EXISTS schema_status;
DROP TYPE IF EXISTS schema_type;
