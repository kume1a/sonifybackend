package playlist

import (
	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
)

type createPlaylistDTO struct {
	Name          string
	ThumbnailPath string
}

type playlistIDDTO struct {
	PlaylistID uuid.UUID
}

type PlaylistFullDTO struct {
	Playlist       sharedmodule.PlaylistDTO
	PlaylistAudios []playlistaudio.PlaylistAudioDTO
}
