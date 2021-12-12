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
	Winner         string `json:"winner"`
	Duration       int64  `json:"duration"`
	DurationVoting int64  `json:"durationVoting"`
}

type listQuizzesResponse struct {
	Page      int32             `json:"page"`
	TotalPage int32             `json:"totalPage"`
	Quizzes   []getQuizResponse `json:"quizzes"`
}

type getAnswerResponse struct {
	ID        string `json:"id"`
	QID       string `json:"qid"`
	Index     int32  `json:"index"`
	Owner     string `json:"owner"`
	Content   string `json:"content"`
	Vote      int32  `json:"vote"`
	Timestamp int64  `json:"timestamp"`
	IsCorrect int32  `json:"isCorrect"`
}

type listAnswersResponse struct {
	Page      int32               `json:"page"`
	TotalPage int32               `json:"totalPage"`
	Answers   []getAnswerResponse `json:"answers"`
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

func toAnswerResponse(answer *db.Answer) getAnswerResponse {
	return getAnswerResponse{
		ID:        answer.ID,
		QID:       answer.QuizID,
		Index:     answer.Index,
		Owner:     answer.Owner,
		Content:   answer.Content.String,
		Vote:      answer.Vote,
		Timestamp: answer.TimestampCreated,
		IsCorrect: answer.IsCorrect,
	}
}

func toListAnswersResponse(page, totalPage int32, answers []db.Answer) listAnswersResponse {
	answersResponse := make([]getAnswerResponse, len(answers))
	for i, v := range answers {
		answersResponse[i] = toAnswerResponse(&v)
	}
	return listAnswersResponse{page, totalPage, answersResponse}
}

func toListQuizzesResponse(page, totalPage int32, quizzes []db.Quiz) listQuizzesResponse {
	quizzesResponse := make([]getQuizResponse, len(quizzes))
	for i, v := range quizzes {
		quizzesResponse[i] = toQuizResponse(&v)
	}
	return listQuizzesResponse{page, totalPage, quizzesResponse}
}
