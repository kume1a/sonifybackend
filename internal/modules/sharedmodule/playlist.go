package sharedmodule

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

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

func PlaylistEntityToDTO(e *database.Playlist) *PlaylistDTO {
	return &PlaylistDTO{
		ID:                e.ID,
		CreatedAt:         e.CreatedAt,
		Name:              e.Name,
		ThumbnailPath:     e.ThumbnailPath.String,
		ThumbnailUrl:      e.ThumbnailUrl.String,
		SpotifyID:         e.SpotifyID.String,
		AudioImportStatus: e.AudioImportStatus,
		AudioCount:        e.AudioCount,
		TotalAudioCount:   e.TotalAudioCount,
	}
}
