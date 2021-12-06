package event

import (
	"context"
	"database/sql"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	store "github.com/tuhalang/quiz-server/app/contracts"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"github.com/tuhalang/quiz-server/app/util"
	"log"
	"math/big"
	"strconv"
)

const (
	chainIdTomoTestNet  = "89"
	contractTomoTestNet = "0x37E807fEEB047C1Fc04Bd1a6B1d35471B2C6bf03"
	tokenTomoTestNet    = "0xBF3Eaf4B19881e00586f279CC946d61142f88fA6"
	eventCreateQuestion = "0x17bc07da768e0d57d22fda497a3a73369d9e20f00932b75abad87b7d7ba1ae12"
	eventPredictAnswer  = "0xae14ab603f57373b14b9e221b39fb04c115320b56157c679eb0c5e74021160f6"
	eventFinish         = "0x45ec5e670f1d367b5dc6f15b95818ac76831881de8f4db9df3d397a7ce8b9331"
)

func (event *QuizEvent) SnapshotPrediction(quizId string, predictionId string, index int) *util.QuizError {
	quizIdByte, err := util.HexTo32Bytes(quizId)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	predictionIdByte, err := util.HexTo32Bytes(predictionId)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	predictionState, quizErr := event.ReadPrediction(*quizIdByte, *predictionIdByte, index)
	if quizErr != nil {
		return quizErr
	}

	prediction, err := event.store.Queries.FindAnswerById(context.Background(), predictionId)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	if prediction.Status == util.StatusDone {
		return nil
	}

	createAnswerParam := db.CreateAnswerParams{
		ID:               predictionId,
		QuizID:           quizId,
		Owner:            predictionState.Owner.String(),
		Content:          sql.NullString{},
		HashContent:      util.ByteToString(predictionState.Answer),
		TimestampCreated: predictionState.Timestamp.Int64(),
		Status:           util.StatusDone,
	}

	if util.Keccak256(prediction.Content.String) == util.ByteToString(predictionState.Answer) {
		createAnswerParam.Content = prediction.Content
	}

	err = event.store.DeleteAnswer(context.Background(), predictionId)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	_, err = event.store.CreateAnswer(context.Background(), createAnswerParam)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	return nil
}

func (event *QuizEvent) SnapshotQuiz(id string) *util.QuizError {
	idBytes, err := util.HexTo32Bytes(id)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}
	quizState, quizErr := event.ReadQuiz(*idBytes)
	if quizErr != nil {
		return quizErr
	}

	quiz, err := event.store.Queries.FindQuizById(context.Background(), id)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	if quiz.Status == util.StatusDone {
		return nil
	}

	createQuizParam := db.CreateQuizParams{
		ID:          id,
		Type:        int32(quizState.QuizType.Int64()),
		Owner:       quizState.Owner.String(),
		Content:     sql.NullString{},
		HashContent: util.ByteToString(quizState.Content),
		Answer:      sql.NullString{},
		HashAnswer: sql.NullString{
			String: util.ByteToString(quizState.Answer),
			Valid:  true,
		},
		Reward: sql.NullInt64{
			Int64: quizState.Reward.Int64(),
			Valid: true,
		},
		Duration: quizState.Duration.Int64(),
		DurationVoting: sql.NullInt64{
			Int64: quizState.DurationVoting.Int64(),
			Valid: true,
		},
		TimestampCreated: quizState.Timestamp.Int64(),
		Status:           util.StatusDone,
	}

	if util.Keccak256(quiz.Content.String) == util.ByteToString(quizState.Content) {
		createQuizParam.Content = quiz.Content
	}
	if util.Keccak256(quiz.Answer.String) == util.ByteToString(quizState.Answer) {
		createQuizParam.Answer = quiz.Answer
	}

	err = event.store.DeleteQuiz(context.Background(), id)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	_, err = event.store.CreateQuiz(context.Background(), createQuizParam)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	return nil
}

func (event *QuizEvent) ReadPrediction(quizId [32]byte, predictionId [32]byte, index int) (*store.QuizGamePrediction, *util.QuizError) {

	instance, quizErr := event.initInstance()
	if quizErr != nil {
		return nil, quizErr
	}

	prediction, err := instance.GetPrediction(&bind.CallOpts{}, quizId, big.NewInt(int64(index)))
	if err != nil {
		return nil, util.NewQuizError(500, err.Error())
	}

	if prediction.Id == predictionId {
		return &prediction, nil
	}

	return nil, util.NewQuizError(404, "prediction not found")
}

func (event *QuizEvent) ReadQuiz(id [32]byte) (*store.QuizGameQuiz, *util.QuizError) {
	instance, quizErr := event.initInstance()
	if quizErr != nil {
		return nil, quizErr
	}

	quiz, err := instance.GetQuiz(&bind.CallOpts{}, id)
	if err != nil {
		return nil, util.NewQuizError(500, err.Error())
	}

	if quiz.Owner.String() == util.AddressZero {
		return nil, util.NewQuizError(404, "Quiz not found")
	}

	return &quiz, nil
}

func (event *QuizEvent) LoadEvents() *util.QuizError {
	events, err := event.Filter(event.config.WssUrl, event.config.ContractAddress, event.config.BlockNumber, event.config.StepNumber)
	if err != nil {
		return util.NewQuizError(500, err.Error())
	}

	if len(events) == 0 {
		err = event.store.UpdateBlockNumber(context.Background(), db.UpdateBlockNumberParams{
			ID:          event.config.ID,
			BlockNumber: event.config.BlockNumber + event.config.StepNumber,
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
			quizErr := event.SnapshotQuiz(id)
			if err != nil {
				return quizErr
			}
		case eventPredictAnswer:
			id := v.Topics[1].String()
			qid := v.Topics[2].String()
			data := util.SplitData(hex.EncodeToString(v.Data))
			indexHex := util.FormatHexNumber(data[0], false)
			index, _ := strconv.ParseInt(indexHex, 16, 64)
			quizErr := event.SnapshotPrediction(qid, id, int(index))
			if err != nil {
				return quizErr
			}
		case eventFinish:
			id := v.Topics[1].String()
			_, err := event.store.FinishQuiz(context.Background(), id)
			if err != nil {
				return util.NewQuizError(500, err.Error())
			}
		}

		err = event.store.UpdateBlockNumber(context.Background(), db.UpdateBlockNumberParams{
			ID:          event.config.ID,
			BlockNumber: int64(v.BlockNumber),
		})
		if err != nil {
			return util.NewQuizError(500, err.Error())
		}
	}

	return nil
}
