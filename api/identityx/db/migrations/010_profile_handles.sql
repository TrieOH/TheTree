-- +goose Up
ALTER TABLE actor_profiles
    ADD COLUMN handle TEXT;

-- unique when present: NULLs (no handle set) never collide
CREATE UNIQUE INDEX uniq_actor_profiles_handle
    ON actor_profiles (handle) WHERE handle IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS uniq_actor_profiles_handle;
ALTER TABLE actor_profiles
    DROP COLUMN IF EXISTS handle;
