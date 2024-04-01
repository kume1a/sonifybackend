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
	DirPublic                    = "./public"
	DirSpotifyAudios             = "./public/spotify_audios"
	DirSpotifyAudioThumbnails    = "./public/spotify_audio_thumbnails"
	DirSpotifyPlaylistThumbnails = "./public/spotify_playlist_thumbnails"
	DirPlaylistThumbnails        = "./public/playlist_thumbnails"
)
