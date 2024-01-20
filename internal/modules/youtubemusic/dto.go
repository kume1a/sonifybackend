package youtubemusic

import (
	"github.com/asaskevich/govalidator"
)

type getYoutubeMusicDto struct {
	VideoID string `json:"videoId" valid:"required"`
}

func (dto *getYoutubeMusicDto) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
