package usersync

import (
	"time"

	"github.com/google/uuid"
)

type userSyncDatumDTO struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `json:"userId"`
	SpotifyLastSyncedAt time.Time `json:"spotifyLastSyncedAt"`
}
