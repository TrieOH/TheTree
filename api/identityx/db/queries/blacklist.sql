-- name: AppendBlacklistEntry :one
INSERT INTO blacklist_entries (created_by_actor_id, project_id, type, target, reason, metadata, expires_at)
VALUES (
   @created_by_actor_id,
   @project_id,
   @type,
   @target,
   @reason,
   @metadata,
   @expires_at
)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetBlacklistEntryByTarget :one
SELECT * FROM blacklist_entries
WHERE target = @target
  AND (expires_at IS NULL OR expires_at > NOW());

-- name: GetBlacklistEntryByTargetAndType :one
SELECT * FROM blacklist_entries
WHERE target = @target
  AND type = @type
  AND (expires_at IS NULL OR expires_at > NOW());

-- name: DeleteExpiredBlacklistEntries :exec
DELETE FROM blacklist_entries WHERE expires_at IS NOT NULL AND expires_at <= NOW();