package youtube

import "github.com/asaskevich/govalidator"

func (dto *downloadYoutubeAudioToUserLibraryDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *downloadYoutubeAudioToPlaylistDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
