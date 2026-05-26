-- +goose Up
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
    CONSTRAINT chk_fields_type CHECK (type IN
        (
            'string',
            'email',
            'int',
            'float',
            'bool',
            'date',
            'time',
            'datetime',
            'select'
        )
    ),

    placeholder JSONB,
    default_value JSONB,

    select_behaviour VARCHAR(8),
    CONSTRAINT chk_fields_select_behaviour CHECK (
        type <> 'select'
        OR select_behaviour IN (
            'checkbox',
            'radio',
            'dropdown-checkbox',
            'dropdown-radio'
        )
    ),

    select_type VARCHAR(8) DEFAULT 'string',
    CONSTRAINT chk_select_type CHECK (
        type <> 'select'
        OR select_type IN (
            'string',
            'email',
            'int',
            'float',
            'bool',
            'date',
            'time',
            'datetime'
        )
    ),

    select_options JSONB,
    CONSTRAINT chk_fields_select_options CHECK (
        type <> 'select'
        OR (
            select_options IS NOT NULL
            AND jsonb_typeof(select_options) = 'array'
        )
    ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_fields_key_format CHECK (key ~ '^[a-z_][a-z0-9_]*$')
);

-- +goose Down
DROP TABLE IF EXISTS fields;