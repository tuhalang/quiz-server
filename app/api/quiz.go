package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"net/http"
)

type getQuizRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *QuizServer) getQuiz(ctx *gin.Context) {
	var req getQuizRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

}

type listQuizzesRequest struct {
	Status int32 `form:"status" binding:"required,min=0,max=1"`
	Page   int32 `form:"page" binding:"required,min=1"`
	Size   int32 `form:"size" biding:"required,min=5,max=30"`
}

func (server *QuizServer) listQuizzes(ctx *gin.Context) {
	var req listQuizzesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
}

type updateQuizRequest struct {
	ID       string `json:"id" binding:"required"`
	Owner    string `json:"owner"`
	Content  string `json:"content" binding:"required"`
	Answer   string `json:"answer"`
	Duration int32  `json:"duration"`
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
		Duration: req.Duration,
	})

	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
	}
	ctx.JSON(http.StatusOK, quiz)
}
