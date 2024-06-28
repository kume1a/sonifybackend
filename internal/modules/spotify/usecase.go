package spotify

import (
	"context"
	"database/sql"
	"log"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/modules/youtube"
	"github.com/kume1a/sonifybackend/internal/shared"

	"github.com/thoas/go-funk"
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

type OnDownloadSpotifyAudioProgress func(progress int, total int)

func DownloadWriteSpotifyAudios(
	ctx context.Context,
	resouceConfig *config.ResourceConfig,
	inputs []DownloadSpotifyAudioInput,
	onProgress OnDownloadSpotifyAudioProgress,
) error {
	spotifyIDs := shared.Map(inputs, func(input DownloadSpotifyAudioInput) string {
		return input.SpotifyID
	})

	dbSpotifyIDs, err := audio.GetAudioSpotifyIdsBySpotifyIds(ctx, resouceConfig.DB, spotifyIDs)
	if err != nil {
		return err
	}

	startingProgress := len(dbSpotifyIDs)
	total := len(inputs)

	if startingProgress == total {
		onProgress(startingProgress, total)
		return nil
	}

	filteredInputs := shared.Where(inputs, func(input DownloadSpotifyAudioInput) bool {
		return !shared.ContainsWhereP(
			dbSpotifyIDs,
			func(dbSpotifyID *database.GetAudioSpotifyIDsBySpotifyIDsRow) bool {
				return dbSpotifyID.SpotifyID.String == input.SpotifyID
			})
	})

	for inputIndex, input := range filteredInputs {
		log.Println("Downloading audio for track: ", input.TrackName, " by ", input.ArtistName, " with Spotify ID: ", input.SpotifyID)

		searchQuery := input.TrackName + " " + input.ArtistName + "\"topic\""

		ytVideoID, err := youtube.GetYoutubeSearchBestMatchVideoID(searchQuery)
		if err != nil {
			return err
		}

		audioOutputPath, _, err := youtube.DownloadYoutubeAudio(ytVideoID, youtube.DownloadYoutubeAudioOptions{
			DownloadThumbnail: false,
		})
		if err != nil {
			return err
		}

		audioFileSize, err := shared.GetFileSize(audioOutputPath)
		if err != nil {
			return err
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

		if err := writeDownloadedSpotifyAudio(ctx, resouceConfig, downloadedSpotifyAudio); err != nil {
			return err
		}

		onProgress(startingProgress+inputIndex+1, total)
	}

	return nil
}

func writeDownloadedSpotifyAudio(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	downloadedSpotifyAudio DownloadedSpotifyAudio,
) error {
	params := database.CreateAudioParams{
		SpotifyID:      sql.NullString{String: downloadedSpotifyAudio.SpotifyID, Valid: true},
		YoutubeVideoID: sql.NullString{String: downloadedSpotifyAudio.YoutubeVideoID, Valid: true},
		Path:           sql.NullString{String: downloadedSpotifyAudio.AudioPath, Valid: true},
		SizeBytes:      sql.NullInt64{Int64: downloadedSpotifyAudio.AudioFileSize.Bytes, Valid: true},
		DurationMs:     sql.NullInt32{Int32: downloadedSpotifyAudio.DurationMs, Valid: true},
		ThumbnailUrl:   sql.NullString{String: downloadedSpotifyAudio.ThumbnailURL, Valid: true},
		Title:          sql.NullString{String: downloadedSpotifyAudio.TrackName, Valid: true},
		Author:         sql.NullString{String: downloadedSpotifyAudio.ArtistName, Valid: true},
	}

	_, err := audio.CreateAudio(ctx, resourceConfig.DB, params)

	return err
}

func mergeSpotifySearchWithDBPlaylists(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	spotifySearch *spotifySearchDTO,
) ([]spotifySearchPlaylistAndDbPlaylist, error) {
	spotifyPlaylistIDs := funk.Map(
		spotifySearch.Playlists.Items,
		func(playlist spotifySearchPlaylistItemDTO) string {
			return playlist.ID
		},
	).([]string)

	dbPlaylists, err := playlist.GetPlaylistsBySpotifyIDs(ctx, resourceConfig.DB, spotifyPlaylistIDs)
	if err != nil {
		return nil, err
	}

	log.Println("Mapping dbPlaylists: ", dbPlaylists)
	return funk.Map(
		spotifySearch.Playlists.Items,
		func(playlist spotifySearchPlaylistItemDTO) spotifySearchPlaylistAndDbPlaylist {
			dbPlaylist := funk.Find(dbPlaylists, func(dbPlaylist database.Playlist) bool {
				return dbPlaylist.SpotifyID.String == playlist.ID
			})

			log.Println("dbPlaylist: ", dbPlaylist)

			if dbPlaylist == nil {
				return spotifySearchPlaylistAndDbPlaylist{
					SpotifySearchPlaylist: playlist,
					DbPlaylist:            nil,
				}
			}

			dbPlaylistValue := dbPlaylist.(database.Playlist)

			return spotifySearchPlaylistAndDbPlaylist{
				SpotifySearchPlaylist: playlist,
				DbPlaylist:            &dbPlaylistValue,
			}
		},
	).([]spotifySearchPlaylistAndDbPlaylist), nil
}
