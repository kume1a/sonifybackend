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
	Title          string    `json:"title"`
	DurationMs     int32     `json:"durationMs"`
	Path           string    `json:"path"`
	Author         string    `json:"author"`
	SizeBytes      int64     `json:"sizeBytes"`
	YoutubeVideoID string    `json:"youtubeVideoId"`
	ThumbnailPath  string    `json:"thumbnailPath"`
	ThumbnailUrl   string    `json:"thumbnailUrl"`
	SpotifyID      string    `json:"spotifyId"`
}

type UserAudioDto struct {
	UserId    uuid.UUID `json:"userId"`
	AudioId   uuid.UUID `json:"audioId"`
	CreatedAt time.Time `json:"createdAt"`
}
