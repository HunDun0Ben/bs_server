package bserr

import "fmt"

type AppError struct {
	Code    int
	Message string
	Data    interface{}
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %d, message: %s, error: %s", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewAppError(code int, message string, data interface{}, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Data:    data,
		Err:     err,
	}
}
