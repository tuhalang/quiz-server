-- name: CreateQuiz :one
INSERT INTO quiz (
    "id",
    "type",
    "owner",
    "content",
    "hash_content",
    "answer",
    "hash_answer",
    "reward",
    "duration",
    "duration_voting",
    "timestamp_created",
    "status"
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: UpdateQuizContent :one
UPDATE quiz SET content = $2 WHERE id = $1 RETURNING *;

-- name: UpdateQuizAnswer :one
UPDATE quiz SET answer = $2 where id = $1 RETURNING *;

-- name: FinishQuiz :one
UPDATE quiz SET status = 0 WHERE id = $1 AND status = 1 RETURNING *;

-- name: FindQuizById :one
SELECT * FROM quiz WHERE id = $1 LIMIT 1;

-- name: CountQuiz :one
SELECT count(*) FROM quiz;

-- name: CountQuizByStatus :one
SELECT count(*) FROM quiz WHERE status = $1;

-- name: FindQuizByStatus :many
SELECT * FROM quiz WHERE status = $1 ORDER BY created_at DESC LIMIT $3 OFFSET $2;

-- name: FindQuizzes :many
SELECT * FROM quiz ORDER BY created_at DESC LIMIT $2 OFFSET $1;

-- name: DeleteQuiz :exec
DELETE FROM quiz WHERE ID = $1;

-- name: UpdateResultQuiz :exec
UPDATE quiz SET status = $2, winner = $3, prediction_winner = $4 where id = $1;