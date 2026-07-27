-- +goose Up
CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    step_id UUID NOT NULL REFERENCES steps(id) ON DELETE CASCADE,
    key VARCHAR(64) NOT NULL,
    CONSTRAINT uniq_key_per_step UNIQUE (step_id, key),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position_hint INT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT false,
    type VARCHAR(16) NOT NULL DEFAULT 'string',
    CONSTRAINT chk_fields_type CHECK (type IN (
        'string', 'email', 'int', 'float', 'bool',
        'date', 'time', 'datetime', 'select', 'file',
        'phone', 'url'
    )),

    placeholder JSONB,
    default_value JSONB,
    config JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_fields_key_format CHECK (key ~ '^[a-z_][a-z0-9_]*$')
);

CREATE TABLE field_select_config (
    field_id UUID PRIMARY KEY REFERENCES fields(id) ON DELETE CASCADE,
    behaviour VARCHAR(32) NOT NULL,
    CONSTRAINT chk_select_behaviour CHECK (behaviour IN (
        'checkbox', 'radio', 'dropdown-checkbox', 'dropdown-radio'
    )),
    value_type VARCHAR(16) NOT NULL DEFAULT 'string',
    CONSTRAINT chk_select_value_type CHECK (value_type IN (
        'string', 'email', 'int', 'float', 'date',
        'time', 'datetime', 'phone', 'url'
    )),
    options JSONB NOT NULL,
    CONSTRAINT chk_select_options CHECK (jsonb_typeof(options) = 'array')
);

-- +goose Down
DROP TABLE IF EXISTS field_select_config;
DROP TABLE IF EXISTS fields;