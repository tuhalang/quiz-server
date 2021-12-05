package service

import (
	"context"
	"database/sql"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"github.com/tuhalang/quiz-server/app/util"
	"log"
)

const (
	// StatusDraft is a status when quiz is unconfirmed
	StatusDraft = -1
)

func (service *QuizService) GetQuiz(id string) (*db.Quiz, *util.QuizError) {
	quiz, err := service.store.Queries.FindQuizById(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, util.NewQuizError(404, "Quiz not found")
		}
		log.Println(err)
		return nil, util.NewQuizError(500, err.Error())
	}
	return &quiz, nil
}

func (service *QuizService) UpdateQuiz(reqQuiz db.Quiz) (*db.Quiz, *util.QuizError) {

	quiz, err := service.store.Queries.FindQuizById(context.Background(), reqQuiz.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			quiz, err = service.store.Queries.CreateQuiz(context.Background(), db.CreateQuizParams{
				ID:          reqQuiz.ID,
				Owner:       reqQuiz.Owner,
				Content:     reqQuiz.Content,
				HashContent: reqQuiz.HashContent,
				Answer:      reqQuiz.Answer,
				HashAnswer:  reqQuiz.HashAnswer,
				Status:      StatusDraft,
			})

			if err != nil {
				return nil, util.NewQuizError(500, err.Error())
			}
			return &quiz, nil
		}
		log.Print(err)
		return nil, util.NewQuizError(500, err.Error())
	}

	if quiz.HashContent == util.Keccak256(reqQuiz.Content.String) {
		quiz, err = service.store.Queries.UpdateQuizContent(context.Background(), db.UpdateQuizContentParams{
			ID:      quiz.ID,
			Content: reqQuiz.Content,
		})
	}

	if reqQuiz.Answer.Valid && quiz.HashAnswer.Valid && quiz.HashAnswer.String == util.Keccak256(reqQuiz.Answer.String) {
		quiz, err = service.store.Queries.UpdateQuizAnswer(context.Background(), db.UpdateQuizAnswerParams{
			ID:     quiz.ID,
			Answer: reqQuiz.Answer,
		})
	}

	return &quiz, nil
}
