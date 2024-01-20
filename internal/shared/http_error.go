package shared

import "net/http"

type HttpError struct {
	Message string
	Code    int
}

func (e HttpError) Error() string {
	return e.Message
}

func ErrBadRequest(msg string) HttpError {
	return HttpError{
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

func StatusUnauthorized(msg string) HttpError {
	return HttpError{
		Message: msg,
		Code:    http.StatusUnauthorized,
	}
}

func StatusForbidden(msg string) HttpError {
	return HttpError{
		Message: msg,
		Code:    http.StatusForbidden,
	}
}

func StatusNotFound(msg string) HttpError {
	return HttpError{
		Message: msg,
		Code:    http.StatusNotFound,
	}
}

func StatusMethodNotAllowed(msg string) HttpError {
	return HttpError{
		Message: msg,
		Code:    http.StatusMethodNotAllowed,
	}
}

func StatusInternalServerError(msg string) HttpError {
	return HttpError{
		Message: msg,
		Code:    http.StatusInternalServerError,
	}
}

func StatusNotImplemented(msg string) HttpError {
	return HttpError{
		Message: msg,
		Code:    http.StatusNotImplemented,
	}
}
