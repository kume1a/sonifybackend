package bgwork

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gocraft/work"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/spotify"
	"github.com/kume1a/sonifybackend/internal/shared"
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
		ctx := context.Background()

		unusedAudios, err := audio.GetUnusedAudios(
			ctx,
			resourceConfig.DB,
		)

		if err != nil {
			log.Println("Error deleting unused audios: ", err)
			return
		}

		for _, unusedAudio := range unusedAudios {
			if err := shared.RunNoResultDBTransaction(
				ctx,
				resourceConfig,
				func(tx *database.Queries) error {
					if err := audio.DeleteAudioByID(
						context.Background(),
						resourceConfig.DB,
						unusedAudio.ID,
					); err != nil {
						log.Println("Error deleting audio: ", err)
						return err
					}

					if err := os.Remove(unusedAudio.Path.String); err != nil {
						log.Println("Error removing unused audio file: ", err)
						return err
					}
					if err := os.Remove(unusedAudio.ThumbnailPath.String); err != nil {
						log.Println("Error removing unused audio thumbnail file: ", err)
						return err
					}

					return nil
				},
			); err != nil {
				log.Println("Error deleting unused audio: ", err)
			}
		}

		deletedCount := len(unusedAudios)
		if deletedCount > 0 {
			log.Println("Deleted ", len(unusedAudios), " unused audios")
		}
	}
}
