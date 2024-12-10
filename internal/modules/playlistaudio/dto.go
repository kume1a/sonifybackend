package playlistaudio

import (
	"github.com/google/uuid"
)

type createPlaylistAudioDTO struct {
	PlaylistID uuid.UUID `json:"playlistId" valid:"required"`
	AudioID    uuid.UUID `json:"audioId" valid:"required"`
}

type deletePlaylistAudioDTO struct {
	PlaylistID uuid.UUID `json:"playlistId" valid:"required"`
	AudioID    uuid.UUID `json:"audioId" valid:"required"`
}
