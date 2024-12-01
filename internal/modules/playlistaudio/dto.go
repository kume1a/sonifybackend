package playlistaudio

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
)

type createPlaylistAudioDTO struct {
	PlaylistID uuid.UUID `json:"playlistId" valid:"required"`
	AudioID    uuid.UUID `json:"audioId" valid:"required"`
}

type deletePlaylistAudioDTO struct {
	PlaylistID uuid.UUID `json:"playlistId" valid:"required"`
	AudioID    uuid.UUID `json:"audioId" valid:"required"`
}

type PlaylistAudioDTO struct {
	ID         uuid.UUID              `json:"id"`
	CreatedAt  time.Time              `json:"createdAt"`
	PlaylistID uuid.UUID              `json:"playlistId"`
	AudioID    uuid.UUID              `json:"audioId"`
	Audio      *sharedmodule.AudioDTO `json:"audio"`
}
