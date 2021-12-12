package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"net/http"
)

func (server *QuizServer) updateAnswer(ctx *gin.Context) {
	var req updateAnswerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	answer, quizErr := server.service.UpdateAnswer(db.Answer{
		ID:     req.ID,
		Index:  req.Index,
		QuizID: req.QID,
		Owner:  "",
		Content: sql.NullString{
			String: req.Content,
			Valid:  true,
		},
		Vote:             0,
		HashContent:      "",
		TimestampCreated: 0,
		Status:           0,
		CreatedAt:        sql.NullTime{},
	})

	server.event.SnapshotPrediction(req.QID, req.ID, int(req.Index))

	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}
	ctx.JSON(http.StatusOK, toAnswerResponse(answer))
}

func (server *QuizServer) listAnswers(ctx *gin.Context) {
	var req listAnswersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	page, totalPage, answers, quizErr := server.service.ListAnswers(req.Page, req.Size, req.QID)
	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}

	ctx.JSON(http.StatusOK, toListAnswersResponse(page, totalPage, answers))
}

func (server *QuizServer) updateVote(ctx *gin.Context) {
	var req updateVoteAnswer
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	quizErr := server.event.SnapshotPrediction(req.QID, req.ID, int(req.Index))

	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}
	ctx.JSON(http.StatusOK, "")
}
