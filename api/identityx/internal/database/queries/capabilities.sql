-- name: CreateCapability :one
INSERT INTO capabilities (project_id, resource, action, created_by)
VALUES (@project_id, @resource, @action, @created_by)
RETURNING *;

-- name: ValidateCapabilities :one
SELECT COUNT(*) = @capability_count AS valid
FROM capabilities
WHERE project_id IS NOT DISTINCT FROM @project_id
    AND id = ANY(@capability_ids::uuid[]);

-- name: ListCapabilitiesByProject :many
SELECT *
FROM capabilities
WHERE project_id IS NOT DISTINCT FROM @project_id;

-- name: AssignCapabilitiesToApiKey :exec
INSERT INTO api_key_capabilities (api_key_id, capability_id, assigned_by, assigned_at)
SELECT @api_key_id, unnest(@capability_ids::uuid[]), @assigned_by, NOW()
ON CONFLICT (api_key_id, capability_id) DO NOTHING;

-- name: ListCapabilityIDsByApiKey :many
SELECT capability_id
FROM api_key_capabilities
WHERE api_key_id = @api_key_id;

-- name: ListCapabilitiesByApiKeyPrefix :many
SELECT c.*
FROM api_key_capabilities akc
JOIN api_keys ak ON ak.id = akc.api_key_id
JOIN capabilities c ON c.id = akc.capability_id
WHERE ak.display_prefix = @display_prefix;

