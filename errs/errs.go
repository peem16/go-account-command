package errs

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFoundError() error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: "Not found",
	}
}

func NewUnexpectedError() error {
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: "InternalServerError",
	}
}

func NewValidationError(message string) error {
	return AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func NewBadRequest(message string) error {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}
