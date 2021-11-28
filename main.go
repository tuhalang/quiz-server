package main

import (
	"database/sql"
	"github.com/tuhalang/quiz-server/app/api"
	"github.com/tuhalang/quiz-server/app/event"
	"github.com/tuhalang/quiz-server/app/service"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DB_SOURCE")
	serverAddr := os.Getenv("SERVER_ADDR")

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	quizEvent, err := event.NewQuizEvent(store)
	if err != nil {
		log.Fatal("cannot init event: ", err)
	}

	quizService, err := service.NewQuizService(store)
	if err != nil {
		log.Fatal("cannot init service: ", err)
	}
	quizServer, err := api.NewQuizServer(quizService, quizEvent)
	if err != nil {
		log.Fatal("cannot init server: ", err)
	}

	err = quizServer.Start(serverAddr)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
