package shared

import (
	"database/sql"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type ApiConfig struct {
	DB    *database.Queries
	SqlDB *sql.DB
}

const (
	DirPublic                    = "public"
	DirYoutubeAudios             = DirPublic + "/youtube_audios"
	DirSpotifyAudioThumbnails    = DirPublic + "/spotify_audio_thumbnails"
	DirSpotifyPlaylistThumbnails = DirPublic + "/spotify_playlist_thumbnails"
	DirPlaylistThumbnails        = DirPublic + "/playlist_thumbnails"
	DirUserLocalAudios           = DirPublic + "/user_local_audios"
	DirUserLocalAudioThumbnails  = DirPublic + "/user_local_audio_thumbnails"
)

func ConfigureGoValidator() {
	govalidator.SetFieldsRequiredByDefault(true)

	// after
	govalidator.CustomTypeTagMap.Set("sliceNotEmpty", func(i interface{}, o interface{}) bool {
		slice, ok := i.([]uuid.UUID)
		if !ok {
			return false
		}
		return len(slice) > 0
	})
}
