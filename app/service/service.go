package service

import db "github.com/tuhalang/quiz-server/app/db/sqlc"

// QuizService is an object service
type QuizService struct {
	store *db.Store
}

// NewQuizService returns a new QuizService
func NewQuizService(store *db.Store) (*QuizService, error) {
	return &QuizService{store: store}, nil
}
