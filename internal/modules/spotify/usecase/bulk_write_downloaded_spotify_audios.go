package spotifyusecase

import (
	"context"
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/database"
	audiousecase "github.com/kume1a/sonifybackend/internal/modules/audio/usecase"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func BulkWriteDownloadedSpotifyAudios(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	downloadedSpotifyAudios []DownloadedSpotifyAudio,
) error {
	params := shared.Map(
		downloadedSpotifyAudios,
		func(downloadedSpotifyAudio DownloadedSpotifyAudio) database.CreateAudioParams {
			return database.CreateAudioParams{
				SpotifyID:      sql.NullString{String: downloadedSpotifyAudio.SpotifyID, Valid: true},
				YoutubeVideoID: sql.NullString{String: downloadedSpotifyAudio.YoutubeVideoID, Valid: true},
				Path:           sql.NullString{String: downloadedSpotifyAudio.AudioPath, Valid: true},
				SizeBytes:      sql.NullInt64{Int64: downloadedSpotifyAudio.AudioFileSize.Bytes, Valid: true},
				DurationMs:     sql.NullInt32{Int32: downloadedSpotifyAudio.DurationMs, Valid: true},
				ThumbnailUrl:   sql.NullString{String: downloadedSpotifyAudio.ThumbnailURL, Valid: true},
				Title:          sql.NullString{String: downloadedSpotifyAudio.TrackName, Valid: true},
				Author:         sql.NullString{String: downloadedSpotifyAudio.ArtistName, Valid: true},
			}
		},
	)

	_, err := audiousecase.BulkWriteAudios(ctx, apiCfg, params)

	return err
}
