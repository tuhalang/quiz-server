package service

import (
	"context"
	"database/sql"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"github.com/tuhalang/quiz-server/app/util"
	"log"
)

func (service *QuizService) UpdateAnswer(reqAnswer db.Answer) (*db.Answer, *util.QuizError) {
	answer, err := service.store.Queries.FindAnswerById(context.Background(), reqAnswer.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			answer, err = service.store.Queries.CreateAnswer(context.Background(), db.CreateAnswerParams{
				ID:               reqAnswer.ID,
				Index:            reqAnswer.Index,
				QuizID:           reqAnswer.QuizID,
				Owner:            "",
				Content:          reqAnswer.Content,
				Vote:             0,
				HashContent:      "",
				TimestampCreated: 0,
				Status:           util.StatusDraft,
			})

			if err != nil {
				return nil, util.NewQuizError(500, err.Error())
			}
			return &answer, nil
		}
		log.Print(err)
		return nil, util.NewQuizError(500, err.Error())
	}

	if answer.HashContent == util.Keccak256(reqAnswer.Content.String) {
		answer, err = service.store.Queries.UpdateAnswerContent(context.Background(), db.UpdateAnswerContentParams{
			ID:      answer.ID,
			Content: reqAnswer.Content,
		})
	}

	return &answer, nil
}
