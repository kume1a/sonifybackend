package auth

import (
	"github.com/asaskevich/govalidator"
)

type googleSignInDTO struct {
	Name     string `json:"name" valid:"stringlength(1|255)"`
	Password string `json:"password" valid:"stringlength(6|255)"`
}

func (dto googleSignInDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
