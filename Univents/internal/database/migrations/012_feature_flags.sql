-- +goose Up

CREATE TABLE feature_definitions (
    key VARCHAR(64) PRIMARY KEY,
    category VARCHAR(32) NOT NULL,
    requires_keys VARCHAR(64)[] NULL,
    config_schema JSONB NULL,
    description TEXT NULL
);

CREATE TABLE feature_flags (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    table_name VARCHAR(64) NOT NULL,
    record_id UUID NOT NULL,
    key VARCHAR(64) NOT NULL REFERENCES feature_definitions(key),
    enabled BOOLEAN NOT NULL DEFAULT false,
    config JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(table_name, record_id, key)
);

CREATE INDEX idx_feature_flags_lookup ON feature_flags(table_name, record_id);

-- Register features
INSERT INTO feature_definitions (key, category, requires_keys, config_schema) VALUES
    ('staff', 'staff', NULL, '{"max_staff": "number"}'),
    ('tokens', 'monetization', NULL, NULL),
    ('staff_scheduler', 'staff', '{staff}'::VARCHAR[], NULL),
    ('paid_activities', 'monetization', '{tokens,activities}'::VARCHAR[], NULL);

-- +goose Down
DROP TABLE IF EXISTS feature_flags;
DROP TABLE IF EXISTS feature_definitions;