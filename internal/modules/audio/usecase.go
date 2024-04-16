package audio

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/youtube"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func BulkWriteAudios(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	params []database.CreateAudioParams,
) ([]database.Audio, error) {
	return shared.RunDbTransaction(
		ctx,
		apiCfg,
		func(tx *database.Queries) ([]database.Audio, error) {
			audios := make([]database.Audio, 0, len(params))

			for _, param := range params {
				audio, err := CreateAudio(ctx, tx, param)
				if err != nil {
					return nil, err
				}

				audios = append(audios, *audio)
			}

			return audios, nil
		},
	)
}

type DownloadYoutubeAudioParams struct {
	ApiConfig *shared.ApiConfig
	Context   context.Context
	UserID    uuid.UUID
	VideoID   string
}

func DownloadYoutubeAudio(params DownloadYoutubeAudioParams) (*database.UserAudio, *database.Audio, *shared.HttpError) {
	// check if the user already has the audio
	if _, err := GetUserAudioByYoutubeVideoId(params.Context, params.ApiConfig.DB, params.UserID, params.VideoID); err == nil {
		return nil, nil, shared.HttpErrConflict(shared.ErrAudioAlreadyExists)
	}

	videoInfo, err := youtube.GetYoutubeVideoInfo(params.VideoID)
	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	filePath, thumbnailPath, err := youtube.DownloadYoutubeAudioWithThumbnail(params.VideoID)
	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	fileSize, err := shared.GetFileSize(filePath)
	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	newAudio, err := CreateAudio(
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
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	userAudio, err := CreateUserAudio(params.Context, params.ApiConfig.DB, database.CreateUserAudioParams{
		UserID:  params.UserID,
		AudioID: newAudio.ID,
	})
	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	return userAudio, newAudio, nil
}
