package shared

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Validatable interface {
	Validate() error
}

func ValidateRequest[T Validatable](r *http.Request) (T, error) {
	var body T

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, errors.New(ErrInvalidJSON)
	}

	if err := body.Validate(); err != nil {
		log.Println(err)
		return body, err
	}

	return body, nil
}
