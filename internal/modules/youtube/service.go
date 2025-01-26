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

type DownloadYoutubeAudioAndSaveToUserLibraryParams struct {
	ApiConfig *config.ApiConfig
	Context   context.Context
	UserID    uuid.UUID
	VideoID   string
}

func DownloadYoutubeAudioAndSaveToUserLibrary(
	params DownloadYoutubeAudioAndSaveToUserLibraryParams,
) (
	*audio.UserAudioWithAudio,
	*shared.HttpError,
) {
	// check if the user already has the audio
	if _, err := useraudio.GetUserAudioByYoutubeVideoID(
		params.Context,
		params.ApiConfig.DB,
		database.GetUserAudioByVideoIDParams{
			UserID:         params.UserID,
			YoutubeVideoID: sql.NullString{String: params.VideoID, Valid: true},
		},
	); err == nil {
		return nil, shared.Conflict(shared.ErrAudioAlreadyExists)
	}

	videoInfo, err := GetYoutubeVideoInfo(params.VideoID)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	filePath, thumbnailPath, err := DownloadYoutubeAudio(
		params.VideoID,
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

	txResult, err := shared.RunDBTransaction(
		params.Context,
		params.ApiConfig.ResourceConfig,
		func(tx *database.Queries) (audio.UserAudioWithAudio, error) {
			newAudio, err := audio.CreateAudio(
				params.Context,
				params.ApiConfig.DB,
				database.CreateAudioParams{
					Title:              sql.NullString{String: strings.TrimSpace(videoInfo.Title), Valid: true},
					Author:             sql.NullString{String: strings.TrimSpace(videoInfo.Uploader), Valid: true},
					DurationMs:         sql.NullInt32{Int32: int32(videoInfo.DurationSeconds * 1000), Valid: true},
					Path:               sql.NullString{String: filePath, Valid: true},
					SizeBytes:          sql.NullInt64{Int64: fileSize.Bytes, Valid: true},
					YoutubeVideoID:     sql.NullString{String: params.VideoID, Valid: true},
					ThumbnailPath:      sql.NullString{String: thumbnailPath, Valid: true},
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

			userAudio, err := useraudio.CreateUserAudio(
				params.Context,
				params.ApiConfig.DB,
				database.CreateUserAudioParams{
					UserID:  params.UserID,
					AudioID: newAudio.ID,
				},
			)
			if err != nil {
				return audio.UserAudioWithAudio{}, shared.InternalServerErrorDef()
			}

			return audio.UserAudioWithAudio{
				UserAudio: userAudio,
				Audio:     newAudio,
			}, nil
		},
	)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return &txResult, nil
}

type DownloadYoutubeAudioAndSaveToPlaylistParams struct {
	ApiConfig  *config.ApiConfig
	Context    context.Context
	UserID     uuid.UUID
	PlaylistID uuid.UUID
	VideoID    string
}

// TODO check if audio already exists, don't just create it in both methods

// TODO attach audio to playlist not user here
func DownloadYoutubeAudioAndSaveToPlaylist(
	params DownloadYoutubeAudioAndSaveToPlaylistParams,
) (
	*playlistaudio.PlaylistAudioWithAudio,
	*shared.HttpError,
) {
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

	// check if the playlist already has the audio
	playlistAudioExists, err := playlistaudio.PlaylistAudioExistsByYoutubeVideoID(
		params.Context,
		params.ApiConfig.DB,
		database.PlaylistAudioExistsByYoutubeVideoIDParams{
			PlaylistID:     params.PlaylistID,
			YoutubeVideoID: sql.NullString{String: params.VideoID, Valid: true},
		},
	)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	if !playlistAudioExists {
		return nil, shared.Conflict(shared.ErrAudioAlreadyExists)
	}

	videoInfo, err := GetYoutubeVideoInfo(params.VideoID)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	filePath, thumbnailPath, err := DownloadYoutubeAudio(
		params.VideoID,
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

	txResult, err := shared.RunDBTransaction(
		params.Context,
		params.ApiConfig.ResourceConfig,
		func(tx *database.Queries) (audio.UserAudioWithAudio, error) {
			newAudio, err := audio.CreateAudio(
				params.Context,
				params.ApiConfig.DB,
				database.CreateAudioParams{
					Title:              sql.NullString{String: strings.TrimSpace(videoInfo.Title), Valid: true},
					Author:             sql.NullString{String: strings.TrimSpace(videoInfo.Uploader), Valid: true},
					DurationMs:         sql.NullInt32{Int32: int32(videoInfo.DurationSeconds * 1000), Valid: true},
					Path:               sql.NullString{String: filePath, Valid: true},
					SizeBytes:          sql.NullInt64{Int64: fileSize.Bytes, Valid: true},
					YoutubeVideoID:     sql.NullString{String: params.VideoID, Valid: true},
					ThumbnailPath:      sql.NullString{String: thumbnailPath, Valid: true},
					PlaylistAudioCount: 1,
					UserAudioCount:     0,
					SpotifyID:          sql.NullString{String: "", Valid: false},
					ThumbnailUrl:       sql.NullString{String: "", Valid: false},
					LocalID:            sql.NullString{String: "", Valid: false},
				},
			)
			if err != nil {
				return audio.UserAudioWithAudio{}, shared.InternalServerErrorDef()
			}

			userAudio, err := useraudio.CreateUserAudio(
				params.Context,
				params.ApiConfig.DB,
				database.CreateUserAudioParams{
					UserID:  params.UserID,
					AudioID: newAudio.ID,
				},
			)
			if err != nil {
				return audio.UserAudioWithAudio{}, shared.InternalServerErrorDef()
			}

			return audio.UserAudioWithAudio{
				UserAudio: userAudio,
				Audio:     newAudio,
			}, nil
		},
	)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return &txResult, nil
}
