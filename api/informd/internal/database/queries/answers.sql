-- name: GetAnswersByResponse :many
SELECT * FROM answers
WHERE response_id = @response_id
ORDER BY answered_at ASC;

-- name: GetAnswersByField :many
SELECT * FROM answers
WHERE field_id = @field_id
ORDER BY answered_at ASC;

-- name: GetAnswersByFormID :many
SELECT a.* FROM answers a
INNER JOIN responses r ON a.response_id = r.id
WHERE r.form_id = @form_id
  AND r.finished_at IS NOT NULL
ORDER BY a.answered_at ASC;

-- name: BatchUpsertAnswers :batchexec
INSERT INTO answers (response_id, field_id, answer, answered_at)
VALUES (@response_id, @field_id, @answer, NOW())
    ON CONFLICT (response_id, field_id) DO UPDATE
        SET answer = EXCLUDED.answer,
            answered_at = NOW();