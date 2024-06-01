package playlist

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type createPlaylistDTO struct {
	Name          string
	ThumbnailPath string
}

type PlaylistDTO struct {
	ID                uuid.UUID              `json:"id"`
	CreatedAt         time.Time              `json:"createdAt"`
	Name              string                 `json:"name"`
	ThumbnailPath     string                 `json:"thumbnailPath"`
	ThumbnailUrl      string                 `json:"thumbnailUrl"`
	SpotifyID         string                 `json:"spotifyId"`
	AudioImportStatus database.ProcessStatus `json:"audioImportStatus"`
	AudioCount        int32                  `json:"audioCount"`
	TotalAudioCount   int32                  `json:"totalAudioCount"`
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
