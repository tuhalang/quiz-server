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

func (service *QuizService) ListAnswers(page int32, size int32, qid string) (int32, int32, []db.Answer, *util.QuizError) {
	totalElements, err := service.store.CountAnswers(context.Background(), qid)
	if err != nil {
		return 0, 0, nil, util.NewQuizError(500, err.Error())
	}

	totalPage := int32(totalElements / int64(size))
	if totalElements%int64(size) != 0 {
		totalPage++
	}

	offset := size * (page - 1)
	answers, err := service.store.FindAnswers(context.Background(), db.FindAnswersParams{
		QuizID: qid,
		Offset: offset,
		Limit:  size,
	})

	if err != nil {
		return 0, 0, nil, util.NewQuizError(500, err.Error())
	}
	return page, totalPage, answers, nil
}

func (service *QuizService) GetCorrectAnswer(quizId string) (*db.Answer, *util.QuizError) {
	answer, err := service.store.GetAnswerCorrect(context.Background(), quizId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, util.NewQuizError(500, err.Error())
	}

	return &answer, nil
}
