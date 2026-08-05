-- name: UpsertEmailTemplate :one
INSERT INTO email_templates (project_id, kind, subject, body)
VALUES (@project_id, @kind, @subject, @body)
ON CONFLICT (project_id, kind)
DO UPDATE SET subject = EXCLUDED.subject, body = EXCLUDED.body, updated_at = NOW()
RETURNING *;

-- name: GetEmailTemplateByProjectAndKind :one
SELECT *
FROM email_templates
WHERE project_id = @project_id
  AND kind = @kind;

-- name: DeleteEmailTemplate :exec
DELETE FROM email_templates
WHERE project_id = @project_id
  AND kind = @kind;

-- name: ListEmailTemplatesByProject :many
SELECT *
FROM email_templates
WHERE project_id = @project_id;
