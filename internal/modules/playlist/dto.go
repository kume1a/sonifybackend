package playlist

import (
	"time"

	"github.com/google/uuid"
)

type createPlaylistDto struct {
	Name          string
	ThumbnailPath string
}

type playlistDto struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Name          string    `json:"name"`
	ThumbnailPath string    `json:"thumbnailPath"`
	ThumbnailUrl  string    `json:"thumbnailUrl"`
	SpotifyId     string    `json:"spotifyId"`
}

type createPlaylistAudioDto struct {
	PlaylistID uuid.UUID `json:"playlistId" valid:"required"`
	AudioID    uuid.UUID `json:"audioId" valid:"required"`
}

type playlistAudioDto struct {
	PlaylistID uuid.UUID `json:"playlistId"`
	AudioID    uuid.UUID `json:"audioId"`
}
