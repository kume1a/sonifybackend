package auth

import "github.com/asaskevich/govalidator"

func (dto *googleSignInDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
