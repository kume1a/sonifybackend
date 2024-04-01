package shared

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

func (dto *KeywordDto) Validate() error {
	if len(dto.Keyword) != 1 {
		return fmt.Errorf("keyword must have exactly one element")
	}

	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *LastCreatedAtPageParamsDto) Validate() error {
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

	accessToken[0] = strings.Replace(accessToken[0], "Bearer ", "", 1)

	return accessToken[0], nil
}
