-- name: CreateAnswer :one
INSERT INTO answer (
    "id",
    "index",
    "quiz_id" ,
    "owner" ,
    "content",
    "vote",
    "hash_content",
    "timestamp_created",
    "status"
) values (
 $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: FindAnswerById :one
SELECT * FROM answer WHERE id = $1 LIMIT 1;

-- name: DeleteAnswer :exec
DELETE FROM answer WHERE ID = $1;

-- name: UpdateAnswerContent :one
UPDATE answer SET content = $2 WHERE id = $1 RETURNING *;

-- name: CountAnswers :one
SELECT count(*) FROM answer WHERE quiz_id = $1;

-- name: FindAnswers :many
SELECT * FROM answer WHERE status = 1 and quiz_id = $1 ORDER BY is_correct desc, vote desc, timestamp_created DESC LIMIT $3 OFFSET $2;

-- name: UpdateAnswerCorrect :exec
UPDATE answer SET is_correct = 1 WHERE id = $1;

-- name: GetAnswerCorrect :one
SELECT * FROM answer WHERE quiz_id = $1 and is_correct = 1 LIMIT 1;

-- name: UpdateVoteNumber :exec
UPDATE answer set vote = $2 where id = $1;