-- +goose Up

ALTER TABLE sessions
    ADD COLUMN user_type TEXT NOT NULL;

ALTER TABLE users
    ADD COLUMN user_type TEXT NOT NULL DEFAULT 'client';

ALTER TABLE sessions
    ADD CONSTRAINT sessions_user_type_check
        CHECK (
            (project_id IS NULL AND user_type = 'client')
                OR
            (project_id IS NOT NULL AND user_type = 'project')
            );

-- +goose Down
ALTER TABLE sessions
DROP CONSTRAINT IF EXISTS sessions_user_type_check;

ALTER TABLE sessions
DROP COLUMN IF EXISTS user_type;

ALTER TABLE users
DROP COLUMN IF EXISTS user_type;
