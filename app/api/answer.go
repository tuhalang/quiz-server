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
	ctx.JSON(http.StatusOK, answer)
}
