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
	AudioIDs []uuid.UUID `json:"audioIds"`
}

type likeAudioDTO struct {
	AudioID uuid.UUID `json:"audioId" valid:"required"`
}

type unlikeAudioDTO struct {
	AudioID uuid.UUID `json:"audioId" valid:"required"`
}

type syncAudioLikesDTO struct {
	AudioIDs []uuid.UUID `json:"audioIds" valid:"required"`
}

type AudioLikeDTO struct {
	UserID  uuid.UUID `json:"userId"`
	AudioID uuid.UUID `json:"audioId"`
}

type AudioDTO struct {
	ID             uuid.UUID     `json:"id"`
	CreatedAt      time.Time     `json:"createdAt"`
	Title          string        `json:"title"`
	DurationMs     int32         `json:"durationMs"`
	Path           string        `json:"path"`
	Author         string        `json:"author"`
	SizeBytes      int64         `json:"sizeBytes"`
	YoutubeVideoID string        `json:"youtubeVideoId"`
	ThumbnailPath  string        `json:"thumbnailPath"`
	ThumbnailUrl   string        `json:"thumbnailUrl"`
	SpotifyID      string        `json:"spotifyId"`
	LocalID        string        `json:"localId"`
	AudioLike      *AudioLikeDTO `json:"likePayload"`
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
