// Code generated by sqlc. DO NOT EDIT.
// source: answer.sql

package db

import (
	"context"
	"database/sql"
)

const countAnswers = `-- name: CountAnswers :one
SELECT count(*) FROM answer WHERE quiz_id = $1
`

func (q *Queries) CountAnswers(ctx context.Context, quizID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAnswers, quizID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAnswer = `-- name: CreateAnswer :one
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
) RETURNING id, index, quiz_id, owner, content, vote, hash_content, timestamp_created, status, is_correct, created_at
`

type CreateAnswerParams struct {
	ID               string         `json:"id"`
	Index            int32          `json:"index"`
	QuizID           string         `json:"quiz_id"`
	Owner            string         `json:"owner"`
	Content          sql.NullString `json:"content"`
	Vote             int32          `json:"vote"`
	HashContent      string         `json:"hash_content"`
	TimestampCreated int64          `json:"timestamp_created"`
	Status           int32          `json:"status"`
}

func (q *Queries) CreateAnswer(ctx context.Context, arg CreateAnswerParams) (Answer, error) {
	row := q.db.QueryRowContext(ctx, createAnswer,
		arg.ID,
		arg.Index,
		arg.QuizID,
		arg.Owner,
		arg.Content,
		arg.Vote,
		arg.HashContent,
		arg.TimestampCreated,
		arg.Status,
	)
	var i Answer
	err := row.Scan(
		&i.ID,
		&i.Index,
		&i.QuizID,
		&i.Owner,
		&i.Content,
		&i.Vote,
		&i.HashContent,
		&i.TimestampCreated,
		&i.Status,
		&i.IsCorrect,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAnswer = `-- name: DeleteAnswer :exec
DELETE FROM answer WHERE ID = $1
`

func (q *Queries) DeleteAnswer(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteAnswer, id)
	return err
}

const findAnswerById = `-- name: FindAnswerById :one
SELECT id, index, quiz_id, owner, content, vote, hash_content, timestamp_created, status, is_correct, created_at FROM answer WHERE id = $1 LIMIT 1
`

func (q *Queries) FindAnswerById(ctx context.Context, id string) (Answer, error) {
	row := q.db.QueryRowContext(ctx, findAnswerById, id)
	var i Answer
	err := row.Scan(
		&i.ID,
		&i.Index,
		&i.QuizID,
		&i.Owner,
		&i.Content,
		&i.Vote,
		&i.HashContent,
		&i.TimestampCreated,
		&i.Status,
		&i.IsCorrect,
		&i.CreatedAt,
	)
	return i, err
}

const findAnswers = `-- name: FindAnswers :many
SELECT id, index, quiz_id, owner, content, vote, hash_content, timestamp_created, status, is_correct, created_at FROM answer WHERE status = 1 and quiz_id = $1 ORDER BY is_correct desc, vote desc, timestamp_created DESC LIMIT $3 OFFSET $2
`

type FindAnswersParams struct {
	QuizID string `json:"quiz_id"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

func (q *Queries) FindAnswers(ctx context.Context, arg FindAnswersParams) ([]Answer, error) {
	rows, err := q.db.QueryContext(ctx, findAnswers, arg.QuizID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Answer
	for rows.Next() {
		var i Answer
		if err := rows.Scan(
			&i.ID,
			&i.Index,
			&i.QuizID,
			&i.Owner,
			&i.Content,
			&i.Vote,
			&i.HashContent,
			&i.TimestampCreated,
			&i.Status,
			&i.IsCorrect,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAnswerCorrect = `-- name: GetAnswerCorrect :one
SELECT id, index, quiz_id, owner, content, vote, hash_content, timestamp_created, status, is_correct, created_at FROM answer WHERE quiz_id = $1 and is_correct = 1 LIMIT 1
`

func (q *Queries) GetAnswerCorrect(ctx context.Context, quizID string) (Answer, error) {
	row := q.db.QueryRowContext(ctx, getAnswerCorrect, quizID)
	var i Answer
	err := row.Scan(
		&i.ID,
		&i.Index,
		&i.QuizID,
		&i.Owner,
		&i.Content,
		&i.Vote,
		&i.HashContent,
		&i.TimestampCreated,
		&i.Status,
		&i.IsCorrect,
		&i.CreatedAt,
	)
	return i, err
}

const updateAnswerContent = `-- name: UpdateAnswerContent :one
UPDATE answer SET content = $2 WHERE id = $1 RETURNING id, index, quiz_id, owner, content, vote, hash_content, timestamp_created, status, is_correct, created_at
`

type UpdateAnswerContentParams struct {
	ID      string         `json:"id"`
	Content sql.NullString `json:"content"`
}

func (q *Queries) UpdateAnswerContent(ctx context.Context, arg UpdateAnswerContentParams) (Answer, error) {
	row := q.db.QueryRowContext(ctx, updateAnswerContent, arg.ID, arg.Content)
	var i Answer
	err := row.Scan(
		&i.ID,
		&i.Index,
		&i.QuizID,
		&i.Owner,
		&i.Content,
		&i.Vote,
		&i.HashContent,
		&i.TimestampCreated,
		&i.Status,
		&i.IsCorrect,
		&i.CreatedAt,
	)
	return i, err
}

const updateAnswerCorrect = `-- name: UpdateAnswerCorrect :exec
UPDATE answer SET is_correct = 1 WHERE id = $1
`

func (q *Queries) UpdateAnswerCorrect(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, updateAnswerCorrect, id)
	return err
}

const updateVoteNumber = `-- name: UpdateVoteNumber :exec
UPDATE answer set vote = $2 where id = $1
`

type UpdateVoteNumberParams struct {
	ID   string `json:"id"`
	Vote int32  `json:"vote"`
}

func (q *Queries) UpdateVoteNumber(ctx context.Context, arg UpdateVoteNumberParams) error {
	_, err := q.db.ExecContext(ctx, updateVoteNumber, arg.ID, arg.Vote)
	return err
}
