package apperror

import (
	"errors"
	"net/http"
)

// AppError represents a domain or application error mapped to an HTTP status code.
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"error"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NotFound(msg string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    msg,
		HTTPStatus: http.StatusNotFound,
	}
}

func Conflict(msg string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    msg,
		HTTPStatus: http.StatusConflict,
	}
}

func Validation(msg string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    msg,
		HTTPStatus: http.StatusBadRequest,
	}
}

func BadRequest(msg string) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    msg,
		HTTPStatus: http.StatusBadRequest,
	}
}

func Internal(msg string, err error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    msg,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

// FromError extracts an *AppError from err, defaulting to 500 Internal Server Error.
func FromError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Internal("internal server error", err)
}
