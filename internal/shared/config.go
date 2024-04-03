package shared

import (
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/database"
)

type ApiConfg struct {
	DB    *database.Queries
	SqlDB *sql.DB
}

const (
	DirPublic                    = "public"
	DirSpotifyAudios             = DirPublic + "/spotify_audios"
	DirSpotifyAudioThumbnails    = DirPublic + "/spotify_audio_thumbnails"
	DirSpotifyPlaylistThumbnails = DirPublic + "/spotify_playlist_thumbnails"
	DirPlaylistThumbnails        = DirPublic + "/playlist_thumbnails"
)
