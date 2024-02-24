package shared

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
)

func (dto *KeywordDto) Validate() error {
	if len(dto.Keyword) != 1 {
		return fmt.Errorf("Keyword must have exactly one element")
	}

	_, err := govalidator.ValidateStruct(dto)
	return err
}

func GetAccessTokenFromRequest(r *http.Request) (string, error) {
	accessToken, ok := r.Header["Authorization"]
	if !ok {
		return "", errors.New(ErrMissingToken)
	}

	if len(accessToken) == 0 {
		return "", errors.New(ErrInvalidToken)
	}

	return accessToken[0], nil
}
