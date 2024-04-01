package playlist

import "github.com/google/uuid"

type createPlaylistDto struct {
	Name          string
	ThumbnailPath string
}

type playlistDto struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ThumbnailPath string    `json:"thumbnailPath"`
}
