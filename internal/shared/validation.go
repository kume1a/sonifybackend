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

func ValidateRequestBody[T Validatable](r *http.Request) (T, error) {
	defer r.Body.Close()

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

func ValidateRequestQuery[T Validatable](r *http.Request) (T, error) {
	var q T

	jsonbody, err := json.Marshal(r.URL.Query())
	if err != nil {
		return q, errors.New(ErrInvalidJSON)
	}

	if err := json.Unmarshal(jsonbody, &q); err != nil {
		return q, errors.New(ErrInvalidJSON)
	}

	if err := q.Validate(); err != nil {
		return q, err
	}

	return q, nil
}
