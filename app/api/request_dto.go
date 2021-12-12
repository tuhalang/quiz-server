package api

type getQuizRequest struct {
	ID string `uri:"id" binding:"required"`
}

type listQuizzesRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" biding:"required,min=5,max=30"`
}

type updateQuizRequest struct {
	ID      string `json:"id" binding:"required"`
	Owner   string `json:"owner"`
	Content string `json:"content" binding:"required"`
	Answer  string `json:"answer"`
}

type updateAnswerRequest struct {
	ID      string `json:"id" binding:"required"`
	QID     string `json:"qid" binding:"required"`
	Index   int32  `json:"index" binding:"min=0"`
	Content string `json:"content" binding:"required"`
}

type updateVoteAnswer struct {
	ID    string `json:"id" binding:"required"`
	QID   string `json:"qid" binding:"required"`
	Index int32  `json:"index" binding:"min=0"`
}

type listAnswersRequest struct {
	QID  string `form:"qid" binding:"required"`
	Page int32  `form:"page" binding:"required,min=1"`
	Size int32  `form:"size" biding:"required,min=5,max=30"`
}
