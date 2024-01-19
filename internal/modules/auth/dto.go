package auth

import (
	"github.com/asaskevich/govalidator"
)

type signInDTO struct {
	Email    string `json:"email" valid:"email,stringlength(1|255)"`
	Password string `json:"password" valid:"stringlength(6|255)"`
}

type signUpDTO struct {
	Email    string `json:"email" valid:"email,stringlength(1|255)"`
	Password string `json:"password" valid:"stringlength(6|255)"`
	Name     string `json:"name" valid:"stringlength(1|255)"`
}

func (dto signInDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto signUpDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
