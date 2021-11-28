-- name: GetEventLog :one
SELECT * FROM event_log
WHERE chain_id = $1 and contract_address = $2 LIMIT 1;

-- name: UpdateBlockNumber :exec
UPDATE event_log
SET block_number = $3
WHERE chain_id = $1 and contract_address = $2;