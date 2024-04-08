package spotifyusecase

import (
	"context"

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

func DownloadSpotifyAudios(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
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
			func(dbSpotifyID *database.GetAudioSpotifyIdsBySpotifyIdsRow) bool {
				return dbSpotifyID.SpotifyID.String == input.SpotifyID
			})
	})

	return shared.ExecuteParallel(
		8,
		filteredInputs,
		func(input DownloadSpotifyAudioInput) (DownloadedSpotifyAudio, error) {
			audioOutputPath, err := shared.NewPublicFileLocation(shared.PublicFileLocationArgs{
				Dir:       shared.DirYoutubeAudios,
				Extension: "mp3",
			})
			if err != nil {
				return DownloadedSpotifyAudio{}, err
			}

			searchQuery := input.TrackName + " " + input.ArtistName

			ytVideoID, err := youtube.GetYoutubeSearchBestMatchVideoID(searchQuery)
			if err != nil {
				return DownloadedSpotifyAudio{}, err
			}

			if err := youtube.DownloadYoutubeAudio(ytVideoID, audioOutputPath); err != nil {
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
