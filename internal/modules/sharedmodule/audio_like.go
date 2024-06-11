package sharedmodule

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type AudioLikeDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    uuid.UUID `json:"userId"`
	AudioID   uuid.UUID `json:"audioId"`
}

func AudioLikeEntityToDTO(e *database.AudioLike) *AudioLikeDTO {
	if e == nil {
		return nil
	}

	return &AudioLikeDTO{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		UserID:    e.UserID,
		AudioID:   e.AudioID,
	}
}

func AudioLikeEntityListToDTOList(e []database.AudioLike) []*AudioLikeDTO {
	res := make([]*AudioLikeDTO, len(e))
	for i, v := range e {
		res[i] = AudioLikeEntityToDTO(&v)
	}

	return res
}
