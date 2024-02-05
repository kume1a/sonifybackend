package youtubemusic

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type getYoutubeMusicDto struct {
	VideoID []string `json:"videoId" valid:"required"`
}

func (dto *getYoutubeMusicDto) Validate() error {
	if len(dto.VideoID) != 1 {
		return fmt.Errorf("VideoID must have exactly one element")
	}

	_, err := govalidator.ValidateStruct(dto)
	return err
}
