package audiolike

import "github.com/asaskevich/govalidator"

func (dto *likeUnlikeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *getAudioLikesDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
