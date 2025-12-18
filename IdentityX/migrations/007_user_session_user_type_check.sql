-- +goose Up

ALTER TABLE user_sessions
    ADD CONSTRAINT user_sessions_user_type_check
        CHECK (
            (project_id IS NULL AND user_type = 'client')
                OR
            (project_id IS NOT NULL AND user_type = 'project')
            );

-- +goose Down
ALTER TABLE user_sessions
DROP CONSTRAINT IF EXISTS user_sessions_user_type_check;