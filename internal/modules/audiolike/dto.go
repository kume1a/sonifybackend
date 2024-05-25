package audiolike

import "github.com/google/uuid"

type likeUnlikeAudioDTO struct {
	AudioID uuid.UUID `json:"audioId" valid:"required"`
}

type getAudioLikesDTO struct {
	IDs uuid.UUIDs `json:"ids" valid:"-"`
}
