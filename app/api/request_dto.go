package api

type getQuizRequest struct {
	ID string `uri:"id" binding:"required"`
}

type listQuizzesRequest struct {
	Status int32 `form:"status" binding:"required,min=0,max=1"`
	Page   int32 `form:"page" binding:"required,min=1"`
	Size   int32 `form:"size" biding:"required,min=5,max=30"`
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
	Index   int32  `json:"index" binding:"required"`
	Content string `json:"content" binding:"required"`
}
