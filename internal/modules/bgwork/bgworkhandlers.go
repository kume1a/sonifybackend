package bgwork

import (
	"context"
	"log"
	"time"

	"github.com/gocraft/work"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/spotify"
)

func CreateHandleDownloadPlaylistAudios(
	resourceConfig *config.ResourceConfig,
) func(job *work.Job) error {
	return func(job *work.Job) error {
		playlistId := job.ArgString("spotifyPlaylistID")
		spotifyAccessToken := job.ArgString("spotifyAccessToken")

		if err := job.ArgError(); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()

		spotify.DownloadSpotifyPlaylistAudios(ctx, resourceConfig, playlistId, spotifyAccessToken)

		return nil
	}
}

func CreateHandleDeleteUnusedAudios(resourceConfig *config.ResourceConfig) func() {
	return func() {
		deletedCount, err := audio.DeleteUnusedAudios(
			context.Background(),
			resourceConfig.DB,
		)

		if err != nil {
			log.Println("Error deleting unused audios: ", err)
			return
		}

		log.Println("Deleted ", deletedCount, " unused audios")
	}
}
