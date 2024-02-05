package shared

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type HttpErrorDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UrlDto struct {
	Url string `json:"url"`
}

type KeywordDto struct {
	Keyword []string `json:"keyword" valid:"required"`
}

func (dto *KeywordDto) Validate() error {
	if len(dto.Keyword) != 1 {
		return fmt.Errorf("Keyword must have exactly one element")
	}

	_, err := govalidator.ValidateStruct(dto)
	return err
}
