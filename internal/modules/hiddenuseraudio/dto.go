package hiddenuseraudio

import (
	"time"

	"github.com/google/uuid"
)

type HiddenUserAudioDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	UserID  uuid.UUID `json:"userId"`
	AudioID uuid.UUID `json:"audioId"`
}
