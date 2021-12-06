-- name: GetChainConfig :one
SELECT * FROM chain_config
WHERE id = $1 LIMIT 1;

-- name: UpdateBlockNumber :exec
UPDATE chain_config
SET block_number = $2
WHERE id = $1;