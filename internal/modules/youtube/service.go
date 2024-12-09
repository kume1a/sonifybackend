package youtube

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/useraudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type DownloadYoutubeAudioParams struct {
	ApiConfig *config.ApiConfig
	Context   context.Context
	UserID    uuid.UUID
	VideoID   string
}

func DownloadYoutubeAudioAndSave(params DownloadYoutubeAudioParams) (
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
