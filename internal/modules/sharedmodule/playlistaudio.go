package sharedmodule

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type PlaylistAudioDTO struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	PlaylistID uuid.UUID `json:"playlistId"`
	AudioID    uuid.UUID `json:"audioId"`
	Audio      *AudioDTO `json:"audio"`
}

type PlaylistAudioWithRelDTO struct {
	*PlaylistAudioDTO
	Audio *AudioDTO `json:"audio"`
}

func PlaylistAudioEntityToDTO(e *database.PlaylistAudio) *PlaylistAudioDTO {
	return &PlaylistAudioDTO{
		ID:         e.ID,
		CreatedAt:  e.CreatedAt,
		PlaylistID: e.PlaylistID,
		AudioID:    e.AudioID,
	}
}
