-- +goose Up
CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,
    namespace_id UUID REFERENCES namespaces(id)
        ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    CONSTRAINT uniq_form_name_per_namespace UNIQUE (namespace_id, name),
    status VARCHAR(8) NOT NULL DEFAULT 'draft',
    CONSTRAINT chk_forms_valid_status CHECK (status IN ('draft', 'open', 'closed', 'archived')),

    opened_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    CONSTRAINT chk_forms_valid_status_state CHECK (
        (status = 'open' AND opened_at IS NOT NULL)
        OR (status = 'closed' AND closed_at IS NOT NULL)
        OR (status = 'archived' AND archived_at IS NOT NULL)
        OR (status = 'draft')
    ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE steps (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    form_id UUID NOT NULL REFERENCES forms(id)
        ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE INDEX idx_forms_owner_id ON forms (owner_id);

CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    step_id UUID NOT NULL REFERENCES steps(id)
        ON DELETE CASCADE,

    key VARCHAR(64) NOT NULL,
    CONSTRAINT uniq_key_per_step UNIQUE (step_id, key),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position_hint INT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT false,

    type VARCHAR(8) NOT NULL DEFAULT 'string',
    CONSTRAINT chk_fields_type CHECK (type IN ('string','email','int','float','bool','select')),
    placeholder JSONB,
    default_value JSONB,

    select_behaviour VARCHAR(8),
    CONSTRAINT chk_fields_select_behaviour CHECK (select_behaviour <> 'select' OR select_behaviour IN ('checkbox', 'radio')),
    select_type VARCHAR(8) NOT NULL DEFAULT 'string',
    CONSTRAINT chk_select_type CHECK (select_type IN ('string','email','int','float','bool','select')),
    select_options JSONB,
    CONSTRAINT chk_fields_select_options CHECK (type <> 'select' OR (select_options IS NOT NULL AND jsonb_typeof(select_options) = 'array')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_fields_key_format CHECK (key ~ '^[a-z_][a-z0-9_]*$')
);

-- +goose Down
DROP TABLE IF EXISTS fields;
DROP INDEX IF EXISTS idx_forms_owner_id;
DROP TABLE IF EXISTS steps;
DROP TABLE IF EXISTS forms;
