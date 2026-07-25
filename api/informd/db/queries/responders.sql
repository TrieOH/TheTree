-- name: CreateResponder :one
INSERT INTO responders (user_id, email)
VALUES (@user_id, @email)
    RETURNING *;

-- name: GetResponderByID :one
SELECT * FROM responders
WHERE id = @id;

-- name: GetResponderByEmail :one
SELECT * FROM responders
WHERE email = @email;

-- name: GetRespondersByFormID :many
SELECT DISTINCT rp.* FROM responders rp
INNER JOIN responses r ON rp.id = r.responder_id
WHERE r.form_id = @form_id
  AND r.finished_at IS NOT NULL;
