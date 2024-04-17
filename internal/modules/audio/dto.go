package audio

import (
	"time"

	"github.com/google/uuid"
)

type downloadYoutubeAudioDTO struct {
	VideoID string `json:"videoId" valid:"required"`
}

type uploadUserLocalMusicDTO struct {
	LocalID       string
	Title         string
	Author        string
	AudioPath     string
	ThumbnailPath string
	DurationMs    int32
}

type getAudiosByIdsDTO struct {
	IDs []uuid.UUID `json:"ids" valid:"required"`
}

type AudioDTO struct {
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
	LocalID        string    `json:"localId"`
}

type UserAudioDTO struct {
	UserId    uuid.UUID `json:"userId"`
	AudioId   uuid.UUID `json:"audioId"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserAudioWithRelDTO struct {
	*UserAudioDTO
	Audio *AudioDTO `json:"audio"`
}
