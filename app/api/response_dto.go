package api

import db "github.com/tuhalang/quiz-server/app/db/sqlc"

type getQuizResponse struct {
	ID             string `json:"id"`
	Content        string `json:"content"`
	Type           int32  `json:"type"`
	Owner          string `json:"owner"`
	Status         int32  `json:"status"`
	Reward         int64  `json:"reward"`
	Timestamp      int64  `json:"timestamp"`
	Winner         string `json:"winner`
	Duration       int64  `json:"duration"`
	DurationVoting int64  `json:"durationVoting"`
}

func toQuizResponse(quiz *db.Quiz) getQuizResponse {
	return getQuizResponse{
		ID:             quiz.ID,
		Content:        quiz.Content.String,
		Type:           quiz.Type,
		Owner:          quiz.Owner,
		Status:         quiz.Status,
		Reward:         quiz.Reward.Int64,
		Timestamp:      quiz.TimestampCreated,
		Winner:         quiz.Winner.String,
		Duration:       quiz.Duration,
		DurationVoting: quiz.DurationVoting.Int64,
	}
}
