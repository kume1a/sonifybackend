package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HttpError struct {
	Message string
	Code    int
}

func (e HttpError) Error() string {
	return e.Message
}

func resJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)

		responseBody := fmt.Sprintf("{message: %s, code: %d}", ErrInternal, http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(responseBody))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func resError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responsing with 500 error, code(%v) %v", code, msg)
	}

	resJson(w, code, HttpErrorDto{
		Message: msg,
		Code:    code,
	})
}

func ResOK(w http.ResponseWriter, payload interface{}) {
	resJson(w, http.StatusOK, payload)
}

func ResCreated(w http.ResponseWriter, payload interface{}) {
	resJson(w, http.StatusCreated, payload)
}

func ResAccepted(w http.ResponseWriter, payload interface{}) {
	resJson(w, http.StatusAccepted, payload)
}

func ResNonAuthoritativeInfo(w http.ResponseWriter, payload interface{}) {
	resJson(w, http.StatusNonAuthoritativeInfo, payload)
}

func ResNoContent(w http.ResponseWriter, payload interface{}) {
	resJson(w, http.StatusNoContent, payload)
}

func ResHttpError(w http.ResponseWriter, httpError HttpError) {
	resError(w, httpError.Code, httpError.Message)
}

func ResBadRequest(w http.ResponseWriter, msg string) {
	resError(w, http.StatusBadRequest, msg)
}

func ResUnauthorized(w http.ResponseWriter, msg string) {
	resError(w, http.StatusUnauthorized, msg)
}

func ResForbidden(w http.ResponseWriter, msg string) {
	resError(w, http.StatusForbidden, msg)
}

func ResNotFound(w http.ResponseWriter, msg string) {
	resError(w, http.StatusNotFound, msg)
}

func ResMethodNotAllowed(w http.ResponseWriter, msg string) {
	resError(w, http.StatusMethodNotAllowed, msg)
}

func ResNotAcceptable(w http.ResponseWriter, msg string) {
	resError(w, http.StatusNotAcceptable, msg)
}

func ResConflict(w http.ResponseWriter, msg string) {
	resError(w, http.StatusConflict, msg)
}

func ResInternalServerError(w http.ResponseWriter, msg string) {
	resError(w, http.StatusInternalServerError, msg)
}

func ResInternalServerErrorDef(w http.ResponseWriter) {
	resError(w, http.StatusInternalServerError, ErrInternal)
}

func ResNotImplemented(w http.ResponseWriter, msg string) {
	resError(w, http.StatusNotImplemented, msg)
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

func HttpErrInternalServerError(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusInternalServerError,
	}
}

func HttpErrInternalServerErrorDef() *HttpError {
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
