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