package shared

import (
	"database/sql"
	"os"

	"github.com/kume1a/sonifybackend/internal/database"
)

type ApiConfg struct {
	DB    *database.Queries
	SqlDB *sql.DB
}

const (
	DirPublic                    = "public"
	DirSpotifyAudios             = DirPublic + string(os.PathSeparator) + "spotify_audios"
	DirSpotifyAudioThumbnails    = DirPublic + string(os.PathSeparator) + "spotify_audio_thumbnails"
	DirSpotifyPlaylistThumbnails = DirPublic + string(os.PathSeparator) + "spotify_playlist_thumbnails"
	DirPlaylistThumbnails        = DirPublic + string(os.PathSeparator) + "playlist_thumbnails"
)
