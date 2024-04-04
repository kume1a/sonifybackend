package usersync

import (
	"github.com/kume1a/sonifybackend/internal/database"
)

func userSyncDatumEntityToDTO(entity *database.UserSyncDatum) userSyncDatumDTO {
	return userSyncDatumDTO{
		ID:                  entity.ID,
		UserID:              entity.UserID,
		SpotifyLastSyncedAt: entity.SpotifyLastSyncedAt,
	}
}
