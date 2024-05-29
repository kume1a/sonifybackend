package audio

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/useraudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type WriteUserImportedLocalMusicParams struct {
	ResourceConfig     *config.ResourceConfig
	Context            context.Context
	UserID             uuid.UUID
	AudioLocalId       string
	AudioTitle         string
	AudioAuthor        string
	AudioPath          string
	AudioThumbnailPath string
	AudioDurationMs    int32
}

func WriteUserImportedLocalMusic(params WriteUserImportedLocalMusicParams) (*UserAudioWithAudio, *shared.HttpError) {
	audioFileSize, err := shared.GetFileSize(params.AudioPath)
	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	res, err := shared.RunDBTransaction(
		params.Context,
		params.ResourceConfig,
		func(tx *database.Queries) (*UserAudioWithAudio, error) {
			audio, err := CreateAudio(
				params.Context,
				tx,
				database.CreateAudioParams{
					Title:         sql.NullString{String: params.AudioTitle, Valid: true},
					Author:        sql.NullString{String: params.AudioAuthor, Valid: params.AudioAuthor != ""},
					Path:          sql.NullString{String: params.AudioPath, Valid: true},
					ThumbnailPath: sql.NullString{String: params.AudioThumbnailPath, Valid: params.AudioThumbnailPath != ""},
					LocalID:       sql.NullString{String: params.AudioLocalId, Valid: params.AudioLocalId != ""},
					DurationMs:    sql.NullInt32{Int32: params.AudioDurationMs, Valid: params.AudioDurationMs != 0},
					SizeBytes:     sql.NullInt64{Int64: audioFileSize.Bytes, Valid: true},
				},
			)
			if err != nil {
				return nil, err
			}

			userAudio, err := useraudio.CreateUserAudio(
				params.Context,
				tx,
				database.CreateUserAudioParams{
					UserID:  params.UserID,
					AudioID: audio.ID,
				},
			)
			if err != nil {
				return nil, err
			}

			return &UserAudioWithAudio{
				UserAudio: userAudio,
				Audio:     audio,
			}, nil
		},
	)

	if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return res, nil
}
