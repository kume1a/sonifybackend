package playlistaudio

import "github.com/asaskevich/govalidator"

func (dto *createPlaylistAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
