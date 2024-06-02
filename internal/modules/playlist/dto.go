package playlist

import (
	"github.com/google/uuid"
)

type createPlaylistDTO struct {
	Name          string
	ThumbnailPath string
}

type playlistIDDTO struct {
	PlaylistID uuid.UUID
}
