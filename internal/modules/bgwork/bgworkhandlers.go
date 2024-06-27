package bgwork

import (
	"context"
	"time"

	"github.com/gocraft/work"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/modules/spotify"
)

func CreateHandleDownloadPlaylistAudios(bgWorkConfig *config.ResourceConfig) func(job *work.Job) error {
	return func(job *work.Job) error {
		playlistId := job.ArgString("spotifyPlaylistID")
		spotifyAccessToken := job.ArgString("spotifyAccessToken")

		if err := job.ArgError(); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()

		spotify.DownloadSpotifyPlaylistAudios(ctx, bgWorkConfig, playlistId, spotifyAccessToken)

		return nil
	}
}
