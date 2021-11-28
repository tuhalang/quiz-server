package event

import (
	"context"
	"database/sql"
	"encoding/hex"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"github.com/tuhalang/quiz-server/app/util"
	"log"
	"os"
	"strconv"
)

const (
	chainIdTomoTestNet  = "89"
	contractTomoTestNet = "0x6EEbaFB7204DA356907671aC8Ab091600d2af9eA"
	tokenTomoTestNet    = "0x15d14572B2505F9fdd47Bb18D5345663ce993672"
	eventCreateQuestion = "0x17bc07da768e0d57d22fda497a3a73369d9e20f00932b75abad87b7d7ba1ae12"
	eventPredictAnswer  = "0xae14ab603f57373b14b9e221b39fb04c115320b56157c679eb0c5e74021160f6"
	eventFinish         = "0x45ec5e670f1d367b5dc6f15b95818ac76831881de8f4db9df3d397a7ce8b9331"
)

func (event *QuizEvent) LoadEvents() *util.QuizError {

	tomoWss := os.Getenv("TOMO_WSS")

	eventLog, err := event.store.GetEventLog(context.Background(), db.GetEventLogParams{
		ChainID:         chainIdTomoTestNet,
		ContractAddress: contractTomoTestNet,
	})
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	events, err := event.Filter(tomoWss, contractTomoTestNet, eventLog.BlockNumber, eventLog.StepNumber)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	if len(events) == 0 {
		err = event.store.UpdateBlockNumber(context.Background(), db.UpdateBlockNumberParams{
			ChainID:         chainIdTomoTestNet,
			ContractAddress: contractTomoTestNet,
			BlockNumber:     eventLog.BlockNumber + eventLog.StepNumber,
		})
		if err != nil {
			return util.NewQuizError(500, err.Error())
		}
	}

	for _, v := range events {
		log.Println("Got a event from topic: ", v.Topics[0].String())
		switch v.Topics[0].String() {
		case eventCreateQuestion:
			id := v.Topics[1].String()
			owner := util.FormatAddress(v.Topics[2].String())
			data := util.SplitData(hex.EncodeToString(v.Data))
			content := data[0]
			answer := data[1]
			timestampHex := util.FormatHexNumber(data[2], false)
			timestamp, _ := strconv.ParseInt(timestampHex, 16, 64)

			createQuizParam := db.CreateQuizParams{
				ID:          id,
				Owner:       owner,
				Content:     sql.NullString{},
				HashContent: content,
				Answer:      sql.NullString{},
				HashAnswer: sql.NullString{
					String: answer,
					Valid:  true,
				},
				TimestampCreated: timestamp,
				Status:           1,
			}

			quiz, err := event.store.Queries.FindQuizById(context.Background(), id)
			if err != nil {
				if err != sql.ErrNoRows {
					return util.NewQuizError(500, err.Error())
				}
			} else {
				if quiz.Status == 1 {
					continue
				}
				if util.Keccak256(quiz.Content.String) == content {
					createQuizParam.Content = quiz.Content
				}
				if util.Keccak256(quiz.Answer.String) == answer {
					createQuizParam.Answer = quiz.Content
				}
				err = event.store.DeleteQuiz(context.Background(), id)
				if err != nil {
					return util.NewQuizError(500, err.Error())
				}
			}

			_, err = event.store.CreateQuiz(context.Background(), createQuizParam)
			if err != nil {
				return util.NewQuizError(500, err.Error())
			}
		case eventPredictAnswer:
			id := v.Topics[1].String()
			qid := v.Topics[2].String()
			owner := v.Topics[3].String()
			data := util.SplitData(hex.EncodeToString(v.Data))
			answer := data[0]
			timestampHex := util.FormatHexNumber(data[1], false)
			timestamp, _ := strconv.ParseInt(timestampHex, 16, 64)

			createAnswerParam := db.CreateAnswerParams{
				ID:               id,
				QuizID:           qid,
				Owner:            owner,
				Content:          sql.NullString{},
				HashContent:      answer,
				TimestampCreated: timestamp,
				Status:           1,
			}

			ans, err := event.store.Queries.FindAnswerById(context.Background(), id)
			if err != nil {
				if err != sql.ErrNoRows {
					return util.NewQuizError(500, err.Error())
				}
			} else {
				if ans.Status == 1 {
					continue
				}
				if util.Keccak256(ans.Content.String) == answer {
					createAnswerParam.Content = ans.Content
				}
				err = event.store.DeleteAnswer(context.Background(), id)
				if err != nil {
					return util.NewQuizError(500, err.Error())
				}
			}

			_, err = event.store.CreateAnswer(context.Background(), createAnswerParam)
			if err != nil {
				return util.NewQuizError(500, err.Error())
			}

		case eventFinish:
			id := v.Topics[1].String()
			_, err := event.store.FinishQuiz(context.Background(), id)
			if err != nil {
				return util.NewQuizError(500, err.Error())
			}
		}

		err = event.store.UpdateBlockNumber(context.Background(), db.UpdateBlockNumberParams{
			ChainID:         chainIdTomoTestNet,
			ContractAddress: contractTomoTestNet,
			BlockNumber:     int64(v.BlockNumber),
		})
		if err != nil {
			return util.NewQuizError(500, err.Error())
		}
	}

	return nil
}
