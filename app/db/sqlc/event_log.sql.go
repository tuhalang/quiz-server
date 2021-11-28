// Code generated by sqlc. DO NOT EDIT.
// source: event_log.sql

package db

import (
	"context"
)

const getEventLog = `-- name: GetEventLog :one
SELECT chain_id, contract_address, block_number, step_number, created_at, updated_at FROM event_log
WHERE chain_id = $1 and contract_address = $2 LIMIT 1
`

type GetEventLogParams struct {
	ChainID         string `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
}

func (q *Queries) GetEventLog(ctx context.Context, arg GetEventLogParams) (EventLog, error) {
	row := q.db.QueryRowContext(ctx, getEventLog, arg.ChainID, arg.ContractAddress)
	var i EventLog
	err := row.Scan(
		&i.ChainID,
		&i.ContractAddress,
		&i.BlockNumber,
		&i.StepNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateBlockNumber = `-- name: UpdateBlockNumber :exec
UPDATE event_log
SET block_number = $3
WHERE chain_id = $1 and contract_address = $2
`

type UpdateBlockNumberParams struct {
	ChainID         string `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	BlockNumber     int64  `json:"block_number"`
}

func (q *Queries) UpdateBlockNumber(ctx context.Context, arg UpdateBlockNumberParams) error {
	_, err := q.db.ExecContext(ctx, updateBlockNumber, arg.ChainID, arg.ContractAddress, arg.BlockNumber)
	return err
}