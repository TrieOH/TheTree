-- +goose Up
CREATE TABLE blacklist_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by_actor_id UUID REFERENCES actors(id)
        ON DELETE SET NULL,
    project_id UUID REFERENCES projects(id)
        ON DELETE CASCADE,

    type TEXT NOT NULL,
    CONSTRAINT chk_blacklist_entries_type CHECK (
        type IN (
            'actor',
            'token',
            'api_key',
            'email',
            'ip'
        )
    ),

    target TEXT NOT NULL,
    reason TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_blacklist_entries_project_id ON blacklist_entries (project_id);
CREATE INDEX idx_blacklist_entries_target ON blacklist_entries (target);
CREATE INDEX idx_blacklist_entries_type ON blacklist_entries (type);
CREATE INDEX idx_blacklist_entries_expires_at ON blacklist_entries (expires_at);
CREATE INDEX idx_blacklist_entries_created_at ON blacklist_entries (created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_blacklist_entries_created_at
DROP INDEX IF EXISTS idx_blacklist_entries_expires_at
DROP INDEX IF EXISTS idx_blacklist_entries_type
DROP INDEX IF EXISTS idx_blacklist_entries_target
DROP INDEX IF EXISTS idx_blacklist_entries_project_id
DROP TABLE IF EXISTS blacklist_entries