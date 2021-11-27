package api

import (
	"github.com/tuhalang/quiz-server/app/service"
	"net/http"
	"os"
)
import "github.com/gin-gonic/gin"

type QuizServer struct {
	service *service.QuizService
	router  *gin.Engine
}

func NewQuizServer(service *service.QuizService) (*QuizServer, error) {
	server := QuizServer{
		service: service,
	}
	server.setupRouter()
	return &server, nil
}

func (server *QuizServer) setupRouter() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	contextPath := os.Getenv("CONTEXT_PATH")
	router.GET(contextPath+"/quizzes/:id", server.getQuiz)
	router.GET(contextPath+"/quizzes", server.listQuizzes)
	router.POST(contextPath+"/quizzes", server.updateQuiz)

	server.router = router
}

func (server *QuizServer) Start(address string) error {
	return server.router.Run(address)
}
