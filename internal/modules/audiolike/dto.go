package audiolike

import "github.com/google/uuid"

type likeUnlikeAudioDTO struct {
	AudioID uuid.UUID `json:"audioId" valid:"required"`
}
