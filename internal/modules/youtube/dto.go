package youtube

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type downloadYoutubeAudioDTO struct {
	VideoID string `json:"videoId" valid:"required"`
}

type getYoutubeMusicUrlDto struct {
	VideoID []string `json:"videoId" valid:"required"`
}

func (dto *getYoutubeMusicUrlDto) Validate() error {
	if len(dto.VideoID) != 1 {
		return fmt.Errorf("VideoID must have exactly one element")
	}

	_, err := govalidator.ValidateStruct(dto)
	return err
}

type youtubeSearchSuggestions struct {
	Query       string   `json:"query"`
	Suggestions []string `json:"suggestions"`
}

type youtubeVideoInfoDTO struct {
	Title           string `json:"title"`
	Uploader        string `json:"uploader"`
	DurationSeconds int    `json:"durationSeconds"`
}
