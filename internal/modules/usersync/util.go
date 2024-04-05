package usersync

import (
	"time"

	"github.com/kume1a/sonifybackend/internal/database"
)

func userSyncDatumEntityToDTO(entity *database.UserSyncDatum) userSyncDatumDTO {
	var spotifyLastSyncedAt *time.Time = nil
	if entity.SpotifyLastSyncedAt.Valid {
		spotifyLastSyncedAt = &entity.SpotifyLastSyncedAt.Time
	}

	return userSyncDatumDTO{
		ID:                  &entity.ID,
		UserID:              &entity.UserID,
		SpotifyLastSyncedAt: spotifyLastSyncedAt,
	}
}
