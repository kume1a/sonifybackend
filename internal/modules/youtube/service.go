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
	*database.UserAudio,
	*database.Audio,
	*shared.HttpError,
) {
	// check if the user already has the audio
	if _, err := useraudio.GetUserAudioByYoutubeVideoId(
		params.Context,
		params.ApiConfig.DB,
		database.GetUserAudioByVideoIDParams{
			UserID:         params.UserID,
			YoutubeVideoID: sql.NullString{String: params.VideoID, Valid: true},
		},
	); err == nil {
		return nil, nil, shared.Conflict(shared.ErrAudioAlreadyExists)
	}

	videoInfo, err := GetYoutubeVideoInfo(params.VideoID)
	if err != nil {
		return nil, nil, shared.InternalServerErrorDef()
	}

	filePath, thumbnailPath, err := DownloadYoutubeAudio(params.VideoID, DownloadYoutubeAudioOptions{
		DownloadThumbnail: true,
	})
	if err != nil {
		return nil, nil, shared.InternalServerErrorDef()
	}

	fileSize, err := shared.GetFileSize(filePath)
	if err != nil {
		return nil, nil, shared.InternalServerErrorDef()
	}

	newAudio, err := audio.CreateAudio(
		params.Context,
		params.ApiConfig.DB,
		database.CreateAudioParams{
			Title:          sql.NullString{String: strings.TrimSpace(videoInfo.Title), Valid: true},
			Author:         sql.NullString{String: strings.TrimSpace(videoInfo.Uploader), Valid: true},
			DurationMs:     sql.NullInt32{Int32: int32(videoInfo.DurationSeconds * 1000), Valid: true},
			Path:           sql.NullString{String: filePath, Valid: true},
			SizeBytes:      sql.NullInt64{Int64: fileSize.Bytes, Valid: true},
			YoutubeVideoID: sql.NullString{String: params.VideoID, Valid: true},
			ThumbnailPath:  sql.NullString{String: thumbnailPath, Valid: true},
		},
	)
	if err != nil {
		return nil, nil, shared.InternalServerErrorDef()
	}

	userAudio, err := useraudio.CreateUserAudio(params.Context, params.ApiConfig.DB, database.CreateUserAudioParams{
		UserID:  params.UserID,
		AudioID: newAudio.ID,
	})
	if err != nil {
		return nil, nil, shared.InternalServerErrorDef()
	}

	return userAudio, newAudio, nil
}
