-- name: CreateQuiz :one
INSERT INTO quiz (
    id,
    owner,
    content,
    hash_content,
    answer,
    hash_answer,
    duration,
    status
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateContent :one
UPDATE quiz SET content = $2 WHERE id = $1 RETURNING *;

-- name: UpdateAnswer :one
UPDATE quiz SET answer = $2 where id = $1 RETURNING *;

-- name: FinishQuiz :one
UPDATE quiz SET status = 0 WHERE id = $1 AND status = 1 RETURNING *;

-- name: FindById :one
SELECT * FROM quiz WHERE id = $1 LIMIT 1;

-- name: Count :one
SELECT count(*) FROM quiz;

-- name: CountByStatus :one
SELECT count(*) FROM quiz WHERE status = $1;

-- name: FindByStatus :many
SELECT * FROM quiz WHERE status = $1 ORDER BY created_at DESC LIMIT $3 OFFSET $2;