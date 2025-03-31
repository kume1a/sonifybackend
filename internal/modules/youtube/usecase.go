package youtube

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/modules/useraudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type downloadYoutubeAudioPayload struct {
	Title          string
	Author         string
	DurationMs     int32
	Path           string
	SizeBytes      int64
	YoutubeVideoID string
	ThumbnailPath  string
}

func downloadYoutubeAudio(youtubeVideoID string) (*downloadYoutubeAudioPayload, error) {
	videoInfo, err := GetYoutubeVideoInfo(youtubeVideoID)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	filePath, thumbnailPath, err := DownloadYoutubeAudio(
		youtubeVideoID,
		DownloadYoutubeAudioOptions{
			DownloadThumbnail: true,
		},
	)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	fileSize, err := shared.GetFileSize(filePath)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return &downloadYoutubeAudioPayload{
		Title:          strings.TrimSpace(videoInfo.Title),
		Author:         strings.TrimSpace(videoInfo.Uploader),
		DurationMs:     int32(videoInfo.DurationSeconds * 1000),
		Path:           filePath,
		SizeBytes:      fileSize.Bytes,
		YoutubeVideoID: youtubeVideoID,
		ThumbnailPath:  thumbnailPath,
	}, nil
}

type DownloadYoutubeAudioAndSaveToUserLibraryParams struct {
	ApiConfig      *config.ApiConfig
	Context        context.Context
	UserID         uuid.UUID
	YoutubeVideoID string
}

func DownloadYoutubeAudioAndSaveToUserLibrary(
	params DownloadYoutubeAudioAndSaveToUserLibraryParams,
) (*audio.UserAudioWithAudio, error) {
	// check if the user already has the audio
	if err := useraudio.ValidateUserAudioDoesNotExist(
		&useraudio.ValidateUserAudioDoesNotExistParams{
			Context:        params.Context,
			ApiConfig:      params.ApiConfig,
			UserID:         params.UserID,
			YoutubeVideoID: params.YoutubeVideoID,
		},
	); err != nil {
		return nil, err
	}

	dbYoutubeAudio, err := audio.GetAudioByYoutubeVideoID(
		params.Context,
		params.ApiConfig.DB,
		params.YoutubeVideoID,
	)

	if err != nil && !shared.IsDBErrorNotFound(err) {
		return nil, shared.InternalServerErrorDef()
	}

	var downloadAudioPayload *downloadYoutubeAudioPayload
	if dbYoutubeAudio == nil || shared.IsDBErrorNotFound(err) {
		downloadAudioPayload, err = downloadYoutubeAudio(params.YoutubeVideoID)
		if err != nil {
			return nil, err
		}
	}

	txResult, err := shared.RunDBTransaction(
		params.Context,
		params.ApiConfig.ResourceConfig,
		func(tx *database.Queries) (audio.UserAudioWithAudio, error) {
			if downloadAudioPayload != nil {
				dbYoutubeAudio, err = audio.CreateAudio(
					params.Context,
					params.ApiConfig.DB,
					database.CreateAudioParams{
						Title:              sql.NullString{String: downloadAudioPayload.Title, Valid: true},
						Author:             sql.NullString{String: downloadAudioPayload.Author, Valid: true},
						DurationMs:         sql.NullInt32{Int32: downloadAudioPayload.DurationMs, Valid: true},
						Path:               sql.NullString{String: downloadAudioPayload.Path, Valid: true},
						SizeBytes:          sql.NullInt64{Int64: downloadAudioPayload.SizeBytes, Valid: true},
						YoutubeVideoID:     sql.NullString{String: downloadAudioPayload.YoutubeVideoID, Valid: true},
						ThumbnailPath:      sql.NullString{String: downloadAudioPayload.ThumbnailPath, Valid: true},
						PlaylistAudioCount: 0,
						UserAudioCount:     1,
						SpotifyID:          sql.NullString{String: "", Valid: false},
						ThumbnailUrl:       sql.NullString{String: "", Valid: false},
						LocalID:            sql.NullString{String: "", Valid: false},
					},
				)
				if err != nil {
					return audio.UserAudioWithAudio{}, shared.InternalServerErrorDef()
				}
			}

			userAudio, err := useraudio.CreateUserAudio(
				params.Context,
				params.ApiConfig.DB,
				database.CreateUserAudioParams{
					UserID:  params.UserID,
					AudioID: dbYoutubeAudio.ID,
				},
			)
			if err != nil {
				return audio.UserAudioWithAudio{}, shared.InternalServerErrorDef()
			}

			return audio.UserAudioWithAudio{
				UserAudio: userAudio,
				Audio:     dbYoutubeAudio,
			}, nil
		},
	)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return &txResult, nil
}

type DownloadYoutubeAudioAndSaveToPlaylistParams struct {
	ApiConfig      *config.ApiConfig
	Context        context.Context
	UserID         uuid.UUID
	PlaylistID     uuid.UUID
	YoutubeVideoID string
}

func DownloadYoutubeAudioAndSaveToPlaylist(
	params DownloadYoutubeAudioAndSaveToPlaylistParams,
) (*playlistaudio.PlaylistAudioWithAudio, error) {
	// check if the playlist already has the audio
	if err := playlistaudio.ValidatePlaylistAudioDoesNotExistByYoutubeVideoID(
		&playlistaudio.ValidatePlaylistAudioDoesNotExistByYoutubeVideoIDParams{
			Context:        params.Context,
			ApiConfig:      params.ApiConfig,
			PlaylistID:     params.PlaylistID,
			YoutubeVideoID: params.YoutubeVideoID,
		},
	); err != nil {
		return nil, err
	}

	// check if playlist belongs to user
	userPlaylistExists, err := sharedmodule.UserPlaylistExists(
		params.Context,
		params.ApiConfig.DB,
		database.UserPlaylistExistsByUserIDAndPlaylistIDParams{
			UserID:     params.UserID,
			PlaylistID: params.PlaylistID,
		},
	)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	if !userPlaylistExists {
		return nil, shared.Forbidden(shared.ErrUserPlaylistNotFound)
	}

	dbYoutubeAudio, err := audio.GetAudioByYoutubeVideoID(
		params.Context,
		params.ApiConfig.DB,
		params.YoutubeVideoID,
	)

	if err != nil && !shared.IsDBErrorNotFound(err) {
		return nil, shared.InternalServerErrorDef()
	}

	var downloadAudioPayload *downloadYoutubeAudioPayload
	if dbYoutubeAudio == nil || shared.IsDBErrorNotFound(err) {
		downloadAudioPayload, err = downloadYoutubeAudio(params.YoutubeVideoID)
		if err != nil {
			return nil, err
		}
	}

	txResult, err := shared.RunDBTransaction(
		params.Context,
		params.ApiConfig.ResourceConfig,
		func(tx *database.Queries) (playlistaudio.PlaylistAudioWithAudio, error) {
			if downloadAudioPayload != nil {
				dbYoutubeAudio, err = audio.CreateAudio(
					params.Context,
					params.ApiConfig.DB,
					database.CreateAudioParams{
						Title:              sql.NullString{String: downloadAudioPayload.Title, Valid: true},
						Author:             sql.NullString{String: downloadAudioPayload.Author, Valid: true},
						DurationMs:         sql.NullInt32{Int32: downloadAudioPayload.DurationMs, Valid: true},
						Path:               sql.NullString{String: downloadAudioPayload.Path, Valid: true},
						SizeBytes:          sql.NullInt64{Int64: downloadAudioPayload.SizeBytes, Valid: true},
						YoutubeVideoID:     sql.NullString{String: downloadAudioPayload.YoutubeVideoID, Valid: true},
						ThumbnailPath:      sql.NullString{String: downloadAudioPayload.ThumbnailPath, Valid: true},
						PlaylistAudioCount: 1,
						UserAudioCount:     0,
						SpotifyID:          sql.NullString{String: "", Valid: false},
						ThumbnailUrl:       sql.NullString{String: "", Valid: false},
						LocalID:            sql.NullString{String: "", Valid: false},
					},
				)
				if err != nil {
					return playlistaudio.PlaylistAudioWithAudio{}, shared.InternalServerErrorDef()
				}
			}

			playlistAudio, err := playlistaudio.CreatePlaylistAudio(
				params.Context,
				params.ApiConfig.DB,
				database.CreatePlaylistAudioParams{
					PlaylistID: params.PlaylistID,
					AudioID:    dbYoutubeAudio.ID,
				},
			)
			if err != nil {
				return playlistaudio.PlaylistAudioWithAudio{}, shared.InternalServerErrorDef()
			}

			return playlistaudio.PlaylistAudioWithAudio{
				PlaylistAudio: playlistAudio,
				Audio:         dbYoutubeAudio,
			}, nil
		},
	)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return &txResult, nil
}
