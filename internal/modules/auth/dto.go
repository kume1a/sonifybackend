package auth

import (
	"github.com/asaskevich/govalidator"
)

type googleSignInDTO struct {
	Token string `json:"token" valid:"required"`
}

func (dto *googleSignInDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
