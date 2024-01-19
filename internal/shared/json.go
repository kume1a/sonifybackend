package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ResJson(w http.ResponseWriter, code int, payload interface{}) {
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

func ResError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responsing with 500 error, code(%v) %v", code, msg)
	}

	ResJson(w, code, HttpErrorDto{
		Message: msg,
		Code:    code,
	})
}
