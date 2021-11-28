-- name: CreateQuiz :one
INSERT INTO quiz (
    id,
    owner,
    content,
    hash_content,
    answer,
    hash_answer,
    timestamp_created,
    status
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
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

-- name: DeleteQuiz :exec
DELETE FROM quiz WHERE ID = $1;