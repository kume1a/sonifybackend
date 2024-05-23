package userplaylist

import "github.com/asaskevich/govalidator"

func (dto *getMyPlaylistsDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
