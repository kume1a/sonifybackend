package userplaylist

import "github.com/asaskevich/govalidator"

func (dto *PlaylistIDsDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
