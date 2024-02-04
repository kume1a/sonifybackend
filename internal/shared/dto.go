package shared

import "github.com/asaskevich/govalidator"

type HttpErrorDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UrlDto struct {
	Url string `json:"url"`
}

type KeywordDto struct {
	Keyword string `json:"keyword" valid:"required"`
}

func (dto *KeywordDto) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
