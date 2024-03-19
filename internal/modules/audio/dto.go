package audio

import (
	"time"

	"github.com/google/uuid"
)

type downloadYoutubeAudioDto struct {
	VideoId string `json:"videoId" valid:"required"`
}

type AudioDto struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Title          string    `json:"title"`
	Duration       int32     `json:"duration"`
	Path           string    `json:"path"`
	Author         string    `json:"author"`
	UserID         uuid.UUID `json:"userId"`
	SizeBytes      int64     `json:"sizeBytes"`
	YoutubeVideoID string    `json:"youtubeVideoId"`
	ThumbnailPath  string    `json:"thumbnailPath"`
}
