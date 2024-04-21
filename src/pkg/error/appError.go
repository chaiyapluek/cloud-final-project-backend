package error

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message ...string) *AppError {
	var errorMessage string
	if len(message) > 0 {
		errorMessage = message[0]
	} else {
		errorMessage = http.StatusText(code)
	}
	return &AppError{
		Code:    code,
		Message: errorMessage,
	}
}

func NewErrBadRequest(msg ...string) *AppError {
	return NewAppError(http.StatusBadRequest, msg...)
}

func NewErrUnauthorized(msg ...string) *AppError {
	return NewAppError(http.StatusUnauthorized, msg...)
}

func NewErrForbidden(msg ...string) *AppError {
	return NewAppError(http.StatusForbidden, msg...)
}

func NewErrNotFound(msg ...string) *AppError {
	return NewAppError(http.StatusNotFound, msg...)
}

func NewErrUnprocessableEntity(msg ...string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, msg...)
}

func NewErrInternalServerError(msg ...string) *AppError {
	return NewAppError(http.StatusInternalServerError, msg...)
}
