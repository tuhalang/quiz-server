package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"net/http"
)

func (server *QuizServer) getQuiz(ctx *gin.Context) {
	var req getQuizRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	quiz, quizErr := server.service.GetQuiz(req.ID)
	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}
	ctx.JSON(http.StatusOK, toQuizResponse(quiz))
}

func (server *QuizServer) updateResultQuiz(ctx *gin.Context) {
	var req getQuizRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	quizErr := server.event.SnapshotQuiz(req.ID)
	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}
	ctx.JSON(http.StatusOK, "")
}

func (server *QuizServer) getWinner(ctx *gin.Context) {
	var req getQuizRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	answer, quizErr := server.service.GetCorrectAnswer(req.ID)
	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}

	if answer == nil {
		ctx.JSON(http.StatusOK, getAnswerResponse{})
		return
	}

	ctx.JSON(http.StatusOK, toAnswerResponse(answer))
}

func (server *QuizServer) listQuizzes(ctx *gin.Context) {
	var req listQuizzesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	page, totalPage, quizzes, quizErr := server.service.ListQuizzes(req.Page, req.Size)
	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}

	ctx.JSON(http.StatusOK, toListQuizzesResponse(page, totalPage, quizzes))
}

func (server *QuizServer) updateQuiz(ctx *gin.Context) {
	var req updateQuizRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	quiz, quizErr := server.service.UpdateQuiz(db.Quiz{
		ID:    req.ID,
		Owner: req.Owner,
		Content: sql.NullString{
			String: req.Content,
			Valid:  true,
		},
		Answer: sql.NullString{
			String: req.Answer,
			Valid:  true,
		},
	})

	server.event.SnapshotQuiz(req.ID)

	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}
	ctx.JSON(http.StatusOK, toQuizResponse(quiz))
}
