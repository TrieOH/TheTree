-- name: GivePlatformRole :one
INSERT INTO platform_roles (actor_id, role, metadata)
VALUES (
    @actor_id,
    @role,
    @metadata
) RETURNING *;
