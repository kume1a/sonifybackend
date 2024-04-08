package spotifyusecase

import (
	"context"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func DownloadWriteSpotifyAudios(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	inputs []DownloadSpotifyAudioInput,
) error {
	downloadedSpotifyAudios, err := DownloadSpotifyAudios(
		ctx,
		apiCfg,
		inputs,
	)
	if err != nil {
		return err
	}

	return BulkWriteDownloadedSpotifyAudios(ctx, apiCfg, downloadedSpotifyAudios)
}
