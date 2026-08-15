-- +goose Up
-- An OAuth identity's scope is its actor's project (platform when NULL).
-- The same provider account may now exist in several scopes: a platform
-- login and a project login with the same Google/GitHub account create
-- separate identities + actors, so the old global UNIQUE (provider, subject)
-- must go. Uniqueness is now per actor — which is per scope.
ALTER TABLE actor_external_identities DROP CONSTRAINT uniq_external_identity;
ALTER TABLE actor_external_identities
    ADD CONSTRAINT uniq_external_identity_per_actor UNIQUE (provider, subject, actor_id);

-- +goose Down
ALTER TABLE actor_external_identities DROP CONSTRAINT uniq_external_identity_per_actor;
ALTER TABLE actor_external_identities
    ADD CONSTRAINT uniq_external_identity UNIQUE (provider, subject);
