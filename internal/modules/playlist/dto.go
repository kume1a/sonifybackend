package playlist

import (
	"time"

	"github.com/google/uuid"
)

type createPlaylistDTO struct {
	Name          string
	ThumbnailPath string
}

type playlistDTO struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Name          string    `json:"name"`
	ThumbnailPath string    `json:"thumbnailPath"`
	ThumbnailUrl  string    `json:"thumbnailUrl"`
	SpotifyId     string    `json:"spotifyId"`
}

type createPlaylistAudioDTO struct {
	PlaylistID uuid.UUID `json:"playlistId" valid:"required"`
	AudioID    uuid.UUID `json:"audioId" valid:"required"`
}

type playlistAudioDTO struct {
	CreatedAt  time.Time `json:"createdAt"`
	PlaylistID uuid.UUID `json:"playlistId"`
	AudioID    uuid.UUID `json:"audioId"`
}

type playlistIDDTO struct {
	PlaylistID uuid.UUID
}
