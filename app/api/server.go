package api

import (
	"github.com/gin-contrib/cors"
	"github.com/tuhalang/quiz-server/app/event"
	"github.com/tuhalang/quiz-server/app/service"
	"net/http"
	"os"
	"time"
)
import "github.com/gin-gonic/gin"

type QuizServer struct {
	service *service.QuizService
	event   *event.QuizEvent
	router  *gin.Engine
}

func NewQuizServer(service *service.QuizService, event *event.QuizEvent) (*QuizServer, error) {
	server := QuizServer{
		service: service,
		event:   event,
	}
	server.setupRouter()
	return &server, nil
}

func (server *QuizServer) setupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	contextPath := os.Getenv("CONTEXT_PATH")
	router.GET(contextPath+"/quizzes/:id", server.getQuiz)
	router.GET(contextPath+"/quizzes", server.listQuizzes)
	router.POST(contextPath+"/quizzes", server.updateQuiz)

	router.GET(contextPath+"/event/:token", server.getEvents)

	server.router = router
}

func (server *QuizServer) Start(address string) error {
	return server.router.Run(address)
}
