package audio

import "github.com/asaskevich/govalidator"

type downloadYoutubeAudioDto struct {
	VideoId string `json:"videoId" valid:"required"`
}

func (dto downloadYoutubeAudioDto) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
