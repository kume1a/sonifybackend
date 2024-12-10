package sharedmodule

import (
	"time"

	"github.com/google/uuid"
)

type PlaylistAudioDTO struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	PlaylistID uuid.UUID `json:"playlistId"`
	AudioID    uuid.UUID `json:"audioId"`
	Audio      *AudioDTO `json:"audio"`
}
