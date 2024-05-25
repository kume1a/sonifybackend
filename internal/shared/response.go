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

type HttpRes struct {
	Code    int
	Payload interface{}
}

func (e HttpError) Error() string {
	return e.Message
}

func ResJson(w http.ResponseWriter, res *HttpRes) {
	data, err := json.Marshal(res.Payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", res.Payload)

		responseBody := fmt.Sprintf("{message: %s, code: %d}", ErrInternal, http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(responseBody))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(res.Code)

	// reflectValue := reflect.ValueOf(res.Payload)
	// if reflectValue.Kind() == reflect.Slice && reflectValue.Len() == 0 {
	// 	w.Write([]byte("[]"))
	// 	return
	// }

	w.Write(data)
}

func ResError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responsing with 500 error, code(%v) %v", code, msg)
	}

	ResJson(
		w,
		&HttpRes{
			Code: code,
			Payload: HttpErrorDto{
				Message: msg,
				Code:    code,
			},
		},
	)
}

func OK(payload interface{}) *HttpRes {
	return &HttpRes{
		Code:    http.StatusOK,
		Payload: payload,
	}
}

func ResOK(w http.ResponseWriter, payload interface{}) {
	ResJson(w, &HttpRes{
		Code:    http.StatusOK,
		Payload: payload,
	})
}

func Created(payload interface{}) *HttpRes {
	return &HttpRes{
		Code:    http.StatusCreated,
		Payload: payload,
	}
}

func ResCreated(w http.ResponseWriter, payload interface{}) {
	ResJson(w, &HttpRes{
		Code:    http.StatusCreated,
		Payload: payload,
	})
}

func Accepted(payload interface{}) *HttpRes {
	return &HttpRes{
		Code:    http.StatusAccepted,
		Payload: payload,
	}
}

func ResAccepted(w http.ResponseWriter, payload interface{}) {
	ResJson(w, &HttpRes{
		Code:    http.StatusAccepted,
		Payload: payload,
	})
}

func NonAuthoritativeInfo(payload interface{}) *HttpRes {
	return &HttpRes{
		Code:    http.StatusNonAuthoritativeInfo,
		Payload: payload,
	}
}

func ResNonAuthoritativeInfo(w http.ResponseWriter, payload interface{}) {
	ResJson(w, &HttpRes{
		Code:    http.StatusNonAuthoritativeInfo,
		Payload: payload,
	})
}

func NoContent() *HttpRes {
	return &HttpRes{
		Code:    http.StatusNoContent,
		Payload: nil,
	}
}

func ResNoContent(w http.ResponseWriter) {
	ResJson(w, &HttpRes{
		Code:    http.StatusNoContent,
		Payload: nil,
	})
}

func ResHttpError(w http.ResponseWriter, httpError *HttpError) {
	ResError(w, httpError.Code, httpError.Message)
}

func ResTryHttpError(w http.ResponseWriter, err error) {
	httpError, ok := err.(*HttpError)
	if !ok {
		ResInternalServerErrorDef(w)
		return
	}

	ResError(w, httpError.Code, httpError.Message)
}

func ResBadRequest(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusBadRequest, msg)
}

func ResUnauthorized(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusUnauthorized, msg)
}

func ResForbidden(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusForbidden, msg)
}

func ResNotFound(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusNotFound, msg)
}

func ResMethodNotAllowed(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusMethodNotAllowed, msg)
}

func ResNotAcceptable(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusNotAcceptable, msg)
}

func ResConflict(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusConflict, msg)
}

func ResInternalServerError(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusInternalServerError, msg)
}

func ResInternalServerErrorDef(w http.ResponseWriter) {
	ResError(w, http.StatusInternalServerError, ErrInternal)
}

func ResNotImplemented(w http.ResponseWriter, msg string) {
	ResError(w, http.StatusNotImplemented, msg)
}

func BadRequest(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

func Unauthorized(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusUnauthorized,
	}
}

func Forbidden(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusForbidden,
	}
}

func NotFound(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotFound,
	}
}

func MethodNotAllowed(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusMethodNotAllowed,
	}
}

func NotAcceptable(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotAcceptable,
	}
}

func Conflict(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusConflict,
	}
}

func InternalServerError(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusInternalServerError,
	}
}

func InternalServerErrorDef() *HttpError {
	return &HttpError{
		Message: ErrInternal,
		Code:    http.StatusInternalServerError,
	}
}

func NotImplemented(msg string) *HttpError {
	return &HttpError{
		Message: msg,
		Code:    http.StatusNotImplemented,
	}
}
