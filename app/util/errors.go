package util

type QuizError struct {
	Message string
	Code    int
}

// NewQuizError returns a new QuizError
func NewQuizError(code int, message string) *QuizError {
	return &QuizError{
		Message: message,
		Code:    code,
	}
}
