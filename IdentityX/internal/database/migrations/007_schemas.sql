-- +goose Up

-- =========================
-- ENUM TYPES
-- =========================

CREATE TYPE schema_type AS ENUM ('core', 'context', 'sub-context');
CREATE TYPE schema_status AS ENUM ('draft', 'published', 'archived');

CREATE TYPE schema_version_status AS ENUM ('draft', 'published', 'archived');

CREATE TYPE field_type AS ENUM (
    'string',
    'email',
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

    based_on_version_id UUID
        REFERENCES schema_versions(id)
            ON DELETE SET NULL,

    UNIQUE (schema_id, version)
);

CREATE INDEX idx_schema_versions_schema_id
    ON schema_versions(schema_id);

CREATE INDEX idx_schema_versions_based_on_version_id
    ON schema_versions(based_on_version_id);

CREATE UNIQUE INDEX uniq_published_schema_versions
    ON schema_versions (schema_id, version)
    WHERE status = 'published';

-- =========================
-- FIELDS
-- =========================

CREATE TABLE schema_fields (
    object_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Stable logical identity of a field within a schema
    id UUID NOT NULL DEFAULT gen_random_uuid(),

    -- (identity scope)
    schema_id UUID NOT NULL
        REFERENCES schemas(id)
            ON DELETE CASCADE,

    schema_version_id UUID NOT NULL
        REFERENCES schema_versions(id)
            ON DELETE CASCADE,

    -- Version-local identifier (can change across versions)
    key VARCHAR(63) NOT NULL,

    type field_type NOT NULL DEFAULT 'string',
    owner field_owner NOT NULL DEFAULT 'system',

    title VARCHAR(255) NOT NULL,
    description TEXT,
    placeholder VARCHAR(255),

    required BOOLEAN NOT NULL DEFAULT false,
    mutable BOOLEAN NOT NULL DEFAULT true,

    default_value JSONB,
    position INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- A field identity appears at most once per version
   UNIQUE (schema_version_id, id),

    -- UI / API guarantees per version
    UNIQUE (schema_version_id, key),
    UNIQUE (schema_version_id, position),

    CHECK (key ~ '^[a-z][a-z0-9_]*$')
);


CREATE INDEX idx_schema_fields_schema_version_id
    ON schema_fields(schema_version_id);

-- =========================
-- FIELD OPTIONS
-- =========================

CREATE TABLE schema_field_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    field_id UUID NOT NULL REFERENCES schema_fields(object_id) ON DELETE CASCADE,

    value VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    position INT NOT NULL,

    UNIQUE (field_id, value)
);

CREATE INDEX idx_schema_field_options_field_id
    ON schema_field_options(field_id);

-- =========================
-- FIELD VISIBILITY RULES
-- =========================

CREATE TABLE schema_field_visibility_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    field_id UUID NOT NULL
        REFERENCES schema_fields(object_id)
            ON DELETE CASCADE,

    depends_on_field_id UUID NOT NULL
        REFERENCES schema_fields(object_id)
            ON DELETE CASCADE,

    operator rule_operator NOT NULL,
    value JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (field_id <> depends_on_field_id)
);

CREATE INDEX idx_visibility_rules_field_id
    ON schema_field_visibility_rules(field_id);

CREATE INDEX idx_visibility_rules_depends_on
    ON schema_field_visibility_rules(depends_on_field_id);

-- =========================
-- FIELD REQUIRED RULES
-- =========================

CREATE TABLE schema_field_required_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    field_id UUID NOT NULL
        REFERENCES schema_fields(object_id)
            ON DELETE CASCADE,

    depends_on_field_id UUID NOT NULL
        REFERENCES schema_fields(object_id)
            ON DELETE CASCADE,

    operator rule_operator NOT NULL,
    value JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (field_id <> depends_on_field_id)
);

CREATE INDEX idx_required_rules_field_id
    ON schema_field_required_rules(field_id);

CREATE INDEX idx_required_rules_depends_on
    ON schema_field_required_rules(depends_on_field_id);

-- =========================
-- CURRENT VERSION FK
-- =========================

ALTER TABLE schemas
    ADD CONSTRAINT fk_current_schema_version
        FOREIGN KEY (current_version_id)
            REFERENCES schema_versions(id)
            ON DELETE SET NULL;

CREATE UNIQUE INDEX one_version_draft_per_schema
    ON schema_versions (schema_id)
    WHERE status = 'draft';


ALTER TABLE sessions
    ADD COLUMN revoked_at TIMESTAMP NULL;

CREATE INDEX sessions_active_idx
    ON sessions (session_id)
    WHERE revoked_at IS NULL;


-- +goose Down
ALTER TABLE DROP COLUMN IF EXISTS revoked_at;
DROP INDEX IF EXISTS sessions_active_idx;
DROP INDEX IF EXISTS idx_schema_versions_based_on_version_id;
DROP INDEX IF EXISTS one_version_draft_per_schema;

ALTER TABLE schemas
DROP CONSTRAINT IF EXISTS fk_current_schema_version;

DROP TABLE IF EXISTS schema_field_required_rules;
DROP TABLE IF EXISTS schema_field_visibility_rules;
DROP TYPE IF EXISTS rule_operator;

DROP TABLE IF EXISTS schema_field_options;

DROP TABLE IF EXISTS schema_fields;
DROP TYPE IF EXISTS field_owner;
DROP TYPE IF EXISTS field_type;

DROP TABLE IF EXISTS schema_versions;
DROP TYPE IF EXISTS schema_version_status;

DROP TABLE IF EXISTS schemas;
DROP TYPE IF EXISTS schema_status;
DROP TYPE IF EXISTS schema_type;