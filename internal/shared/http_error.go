package shared

import (
	"net/http"
)

type HttpError struct {
	Message string
	Code    int
}

func (e HttpError) Error() string {
	return e.Message
}

func HttpErrBadRequest(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

func HttpErrUnauthorized(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusUnauthorized,
	}
}

func HttpErrForbidden(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusForbidden,
	}
}

func HttpErrNotFound(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotFound,
	}
}

func HttpErrMethodNotAllowed(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusMethodNotAllowed,
	}
}

func HttpErrInternalServerError() *HttpError {
	return &HttpError{
		Message: ErrInternal,
		Code:    http.StatusInternalServerError,
	}
}

func HttpErrNotImplemented(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotImplemented,
	}
}
