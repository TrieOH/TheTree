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

-- name: PutBlacklistEntry :exec
-- Adds the entry, or restarts it when it is already there: created_at moves
-- to now, so an actor entry blocks every token signed before this call.
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
ON CONFLICT (type, target, project_id) DO UPDATE
SET created_at = NOW(), reason = EXCLUDED.reason, expires_at = EXCLUDED.expires_at;

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