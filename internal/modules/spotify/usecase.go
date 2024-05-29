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
	resourceConfig *config.ResourceConfig,
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

	_, err := audio.BulkCreateAudios(ctx, resourceConfig, params)

	return err
}

func DownloadSpotifyAudios(
	ctx context.Context,
	resouceConfig *config.ResourceConfig,
	inputs []DownloadSpotifyAudioInput,
) ([]DownloadedSpotifyAudio, error) {
	spotifyIDs := shared.Map(inputs, func(input DownloadSpotifyAudioInput) string {
		return input.SpotifyID
	})

	dbSpotifyIDs, err := audio.GetAudioSpotifyIdsBySpotifyIds(ctx, resouceConfig.DB, spotifyIDs)
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

	downloadedSpotifyAudios := []DownloadedSpotifyAudio{}
	for _, input := range filteredInputs {
		searchQuery := input.TrackName + " " + input.ArtistName + "\"topic\""

		ytVideoID, err := youtube.GetYoutubeSearchBestMatchVideoID(searchQuery)
		if err != nil {
			return []DownloadedSpotifyAudio{}, err
		}

		audioOutputPath, _, err := youtube.DownloadYoutubeAudio(ytVideoID, youtube.DownloadYoutubeAudioOptions{
			DownloadThumbnail: false,
		})
		if err != nil {
			return []DownloadedSpotifyAudio{}, err
		}

		audioFileSize, err := shared.GetFileSize(audioOutputPath)
		if err != nil {
			return []DownloadedSpotifyAudio{}, err
		}

		downloadedSpotifyAudio := DownloadedSpotifyAudio{
			AudioPath:      audioOutputPath,
			AudioFileSize:  audioFileSize,
			YoutubeVideoID: ytVideoID,
			SpotifyID:      input.SpotifyID,
			DurationMs:     input.DurationMs,
			ThumbnailURL:   input.ThumbnailURL,
			TrackName:      input.TrackName,
			ArtistName:     input.ArtistName,
		}

		downloadedSpotifyAudios = append(downloadedSpotifyAudios, downloadedSpotifyAudio)
	}

	return downloadedSpotifyAudios, nil
}

func DownloadWriteSpotifyAudios(
	ctx context.Context,
	resouceConfig *config.ResourceConfig,
	inputs []DownloadSpotifyAudioInput,
) error {
	downloadedSpotifyAudios, err := DownloadSpotifyAudios(
		ctx,
		resouceConfig,
		inputs,
	)
	if err != nil {
		return err
	}

	return BulkWriteDownloadedSpotifyAudios(ctx, resouceConfig, downloadedSpotifyAudios)
}
