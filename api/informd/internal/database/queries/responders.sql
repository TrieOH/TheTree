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
