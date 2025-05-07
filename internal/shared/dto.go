package shared

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type HttpErrorDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UrlDTO struct {
	Url string `json:"url"`
}

type KeywordDTO struct {
	Keyword []string `json:"keyword" valid:"required"`
}

type AudioIDDTO struct {
	AudioID string `json:"audioId" valid:"required,uuid"`
}

type AudioIDsDTO struct {
	AudioIDs []uuid.UUID `json:"audioIds" valid:"-"`
}

type OkDTO struct {
	Ok bool `json:"ok"`
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

type OptionalIDsDTO struct {
	IDs uuid.UUIDs `json:"ids" valid:"-"`
}

// TODO fix validation
type RequiredIDsDTO struct {
	IDs uuid.UUIDs `json:"ids" valid:"-"`
}

func (dto *KeywordDTO) Validate() error {
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

func (dto *OptionalIDsDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *RequiredIDsDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *AudioIDDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *AudioIDsDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
