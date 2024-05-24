package audiolike

import "github.com/google/uuid"

type AudioLikeDTO struct {
	UserID  uuid.UUID `json:"userId"`
	AudioID uuid.UUID `json:"audioId"`
}

type likeAudioDTO struct {
	AudioID uuid.UUID `json:"audioId" valid:"required"`
}

type unlikeAudioDTO struct {
	AudioID uuid.UUID `json:"audioId" valid:"required"`
}

type getAudioLikesDTO struct {
	IDs uuid.UUIDs `json:"ids" valid:"-"`
}
