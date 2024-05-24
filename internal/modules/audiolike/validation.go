package audiolike

import "github.com/asaskevich/govalidator"

func (dto *likeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *unlikeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *getAudioLikesDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
