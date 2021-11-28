package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *QuizServer) getEvents(ctx *gin.Context) {
	quizErr := server.event.LoadEvents()
	if quizErr != nil {
		ctx.JSON(quizErr.Code, quizErr)
		return
	}
	ctx.JSON(http.StatusOK, "Success")
}
