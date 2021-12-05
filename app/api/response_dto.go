package api

import db "github.com/tuhalang/quiz-server/app/db/sqlc"

type getQuizResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func toQuizResponse(quiz *db.Quiz) getQuizResponse {
	return getQuizResponse{
		ID:      quiz.ID,
		Content: quiz.Content.String,
	}
}
