// Code generated by sqlc. DO NOT EDIT.
// source: quiz.sql

package db

import (
	"context"
	"database/sql"
)

const count = `-- name: Count :one
SELECT count(*) FROM quiz
`

func (q *Queries) Count(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, count)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countByStatus = `-- name: CountByStatus :one
SELECT count(*) FROM quiz WHERE status = $1
`

func (q *Queries) CountByStatus(ctx context.Context, status int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, countByStatus, status)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createQuiz = `-- name: CreateQuiz :one
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
) RETURNING id, owner, content, hash_content, answer, hash_answer, duration, status, created_at
`

type CreateQuizParams struct {
	ID          string         `json:"id"`
	Owner       string         `json:"owner"`
	Content     sql.NullString `json:"content"`
	HashContent string         `json:"hash_content"`
	Answer      sql.NullString `json:"answer"`
	HashAnswer  sql.NullString `json:"hash_answer"`
	Duration    int32          `json:"duration"`
	Status      int32          `json:"status"`
}

func (q *Queries) CreateQuiz(ctx context.Context, arg CreateQuizParams) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, createQuiz,
		arg.ID,
		arg.Owner,
		arg.Content,
		arg.HashContent,
		arg.Answer,
		arg.HashAnswer,
		arg.Duration,
		arg.Status,
	)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Content,
		&i.HashContent,
		&i.Answer,
		&i.HashAnswer,
		&i.Duration,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const findById = `-- name: FindById :one
SELECT id, owner, content, hash_content, answer, hash_answer, duration, status, created_at FROM quiz WHERE id = $1 LIMIT 1
`

func (q *Queries) FindById(ctx context.Context, id string) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, findById, id)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Content,
		&i.HashContent,
		&i.Answer,
		&i.HashAnswer,
		&i.Duration,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const findByStatus = `-- name: FindByStatus :many
SELECT id, owner, content, hash_content, answer, hash_answer, duration, status, created_at FROM quiz WHERE status = $1 ORDER BY created_at DESC LIMIT $3 OFFSET $2
`

type FindByStatusParams struct {
	Status int32 `json:"status"`
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) FindByStatus(ctx context.Context, arg FindByStatusParams) ([]Quiz, error) {
	rows, err := q.db.QueryContext(ctx, findByStatus, arg.Status, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Quiz
	for rows.Next() {
		var i Quiz
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Content,
			&i.HashContent,
			&i.Answer,
			&i.HashAnswer,
			&i.Duration,
			&i.Status,
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

const finishQuiz = `-- name: FinishQuiz :one
UPDATE quiz SET status = 0 WHERE id = $1 AND status = 1 RETURNING id, owner, content, hash_content, answer, hash_answer, duration, status, created_at
`

func (q *Queries) FinishQuiz(ctx context.Context, id string) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, finishQuiz, id)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Content,
		&i.HashContent,
		&i.Answer,
		&i.HashAnswer,
		&i.Duration,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const updateAnswer = `-- name: UpdateAnswer :one
UPDATE quiz SET answer = $2 where id = $1 RETURNING id, owner, content, hash_content, answer, hash_answer, duration, status, created_at
`

type UpdateAnswerParams struct {
	ID     string         `json:"id"`
	Answer sql.NullString `json:"answer"`
}

func (q *Queries) UpdateAnswer(ctx context.Context, arg UpdateAnswerParams) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, updateAnswer, arg.ID, arg.Answer)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Content,
		&i.HashContent,
		&i.Answer,
		&i.HashAnswer,
		&i.Duration,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const updateContent = `-- name: UpdateContent :one
UPDATE quiz SET content = $2 WHERE id = $1 RETURNING id, owner, content, hash_content, answer, hash_answer, duration, status, created_at
`

type UpdateContentParams struct {
	ID      string         `json:"id"`
	Content sql.NullString `json:"content"`
}

func (q *Queries) UpdateContent(ctx context.Context, arg UpdateContentParams) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, updateContent, arg.ID, arg.Content)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Content,
		&i.HashContent,
		&i.Answer,
		&i.HashAnswer,
		&i.Duration,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}