package playlist

import (
	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
)

type playlistIDDTO struct {
	PlaylistID uuid.UUID
}

type PlaylistFullDTO struct {
	sharedmodule.PlaylistDTO
	PlaylistAudios []playlistaudio.PlaylistAudioDTO `json:"playlistAudios"`
}
