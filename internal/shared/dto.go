package shared

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
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

type UserDto struct {
	ID           uuid.UUID             `json:"id"`
	CreatedAt    time.Time             `json:"createdAt"`
	Name         string                `json:"name"`
	Email        string                `json:"email"`
	AuthProvider database.AuthProvider `json:"authProvider"`
}

type LastCreatedAtPageParamsDto struct {
	LastCreatedAt time.Time `json:"lastCreatedAt"`
	Limit         int32     `json:"limit" valid:"required,max(200)"`
}

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
