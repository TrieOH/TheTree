-- name: GetAnswersByResponse :many
SELECT * FROM answers
WHERE response_id = @response_id
ORDER BY answered_at ASC;

-- name: BatchUpsertAnswers :batchexec
INSERT INTO answers (response_id, field_id, answer, answered_at)
VALUES (@response_id, @field_id, @answer, NOW())
    ON CONFLICT (response_id, field_id) DO UPDATE
        SET answer = EXCLUDED.answer,
            answered_at = NOW();