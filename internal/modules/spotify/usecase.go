package spotify

import (
	"context"
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/youtube"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type DownloadSpotifyAudioInput struct {
	TrackName    string
	ArtistName   string
	SpotifyID    string
	DurationMs   int32
	ThumbnailURL string
}

type DownloadedSpotifyAudio struct {
	AudioPath     string
	AudioFileSize *shared.FileSize

	YoutubeVideoID string
	SpotifyID      string
	DurationMs     int32
	ThumbnailURL   string
	TrackName      string
	ArtistName     string
}

func BulkWriteDownloadedSpotifyAudios(
	ctx context.Context,
	apiCfg *config.ApiConfig,
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

	_, err := audio.BulkCreateAudios(ctx, apiCfg, params)

	return err
}

func DownloadSpotifyAudios(
	ctx context.Context,
	apiCfg *config.ApiConfig,
	inputs []DownloadSpotifyAudioInput,
) ([]DownloadedSpotifyAudio, error) {
	spotifyIDs := shared.Map(inputs, func(input DownloadSpotifyAudioInput) string {
		return input.SpotifyID
	})

	dbSpotifyIDs, err := audio.GetAudioSpotifyIdsBySpotifyIds(ctx, apiCfg.DB, spotifyIDs)
	if err != nil {
		return nil, err
	}

	filteredInputs := shared.Where(inputs, func(input DownloadSpotifyAudioInput) bool {
		return !shared.ContainsWhereP(
			dbSpotifyIDs,
			func(dbSpotifyID *database.GetAudioSpotifyIDsBySpotifyIDsRow) bool {
				return dbSpotifyID.SpotifyID.String == input.SpotifyID
			})
	})

	return shared.ExecuteParallel(
		3,
		filteredInputs,
		func(input DownloadSpotifyAudioInput) (DownloadedSpotifyAudio, error) {
			searchQuery := input.TrackName + " " + input.ArtistName + "\"topic\""

			ytVideoID, err := youtube.GetYoutubeSearchBestMatchVideoID(searchQuery)
			if err != nil {
				return DownloadedSpotifyAudio{}, err
			}

			audioOutputPath, _, err := youtube.DownloadYoutubeAudio(ytVideoID, youtube.DownloadYoutubeAudioOptions{
				DownloadThumbnail: false,
			})
			if err != nil {
				return DownloadedSpotifyAudio{}, err
			}

			audioFileSize, err := shared.GetFileSize(audioOutputPath)
			if err != nil {
				return DownloadedSpotifyAudio{}, err
			}

			return DownloadedSpotifyAudio{
				AudioPath:      audioOutputPath,
				AudioFileSize:  audioFileSize,
				YoutubeVideoID: ytVideoID,
				SpotifyID:      input.SpotifyID,
				DurationMs:     input.DurationMs,
				ThumbnailURL:   input.ThumbnailURL,
				TrackName:      input.TrackName,
				ArtistName:     input.ArtistName,
			}, nil
		},
	)
}

func DownloadWriteSpotifyAudios(
	ctx context.Context,
	apiCfg *config.ApiConfig,
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
