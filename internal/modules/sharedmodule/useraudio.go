package sharedmodule

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type UserAudioDTO struct {
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"userId"`
	AudioId   uuid.UUID `json:"audioId"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserAudioWithRelDTO struct {
	*UserAudioDTO
	Audio *AudioDTO `json:"audio"`
}

func UserAudioEntityToDTO(e *database.UserAudio) *UserAudioDTO {
	return &UserAudioDTO{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		UserId:    e.UserID,
		AudioId:   e.AudioID,
	}
}

func UserAudioEntitiesToDTOs(entities []database.UserAudio) []*UserAudioDTO {
	dtos := make([]*UserAudioDTO, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, UserAudioEntityToDTO(&entity))
	}
	return dtos
}
