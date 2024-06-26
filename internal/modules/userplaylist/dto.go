package userplaylist

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
)

type userPlaylistDTO struct {
	ID                     uuid.UUID                 `json:"id"`
	UserID                 uuid.UUID                 `json:"userId"`
	PlaylistID             uuid.UUID                 `json:"playlistId"`
	CreatedAt              time.Time                 `json:"createdAt"`
	IsSpotifySavedPlaylist bool                      `json:"isSpotifySavedPlaylist"`
	Playlist               *sharedmodule.PlaylistDTO `json:"playlist"`
}

type PlaylistIDsDTO struct {
	PlaylistIDs []uuid.UUID `json:"playlistIds" valid:"-"`
}
