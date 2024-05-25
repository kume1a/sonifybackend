package youtube

import "github.com/asaskevich/govalidator"

func (dto *downloadYoutubeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
