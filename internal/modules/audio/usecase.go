package audio

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
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
					Title:              sql.NullString{String: params.AudioTitle, Valid: true},
					Author:             sql.NullString{String: params.AudioAuthor, Valid: params.AudioAuthor != ""},
					Path:               sql.NullString{String: params.AudioPath, Valid: true},
					ThumbnailPath:      sql.NullString{String: params.AudioThumbnailPath, Valid: params.AudioThumbnailPath != ""},
					LocalID:            sql.NullString{String: params.AudioLocalId, Valid: params.AudioLocalId != ""},
					DurationMs:         sql.NullInt32{Int32: params.AudioDurationMs, Valid: params.AudioDurationMs != 0},
					SizeBytes:          sql.NullInt64{Int64: audioFileSize.Bytes, Valid: true},
					PlaylistAudioCount: 0,
					UserAudioCount:     1,
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

func WriteInitialAudioRelCount(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
) *shared.HttpError {
	if _, err := shared.RunDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) (*UserAudioWithAudio, error) {
			audioIDs, err := GetAllAudioIDs(ctx, tx)
			if err != nil {
				return nil, err
			}

			for _, audioID := range audioIDs {
				userAudioCount, err := useraudio.CountUserAudiosByAudioID(ctx, tx, audioID)
				if err != nil {
					return nil, err
				}

				playlistAudioCount, err := playlistaudio.CountPlaylistAudiosByAudioID(ctx, tx, audioID)
				if err != nil {
					return nil, err
				}

				if _, err := UpdateAudioByID(ctx, tx, database.UpdateAudioByIDParams{
					AudioID:            uuid.NullUUID{UUID: audioID, Valid: true},
					PlaylistAudioCount: sql.NullInt32{Int32: int32(playlistAudioCount), Valid: true},
					UserAudioCount:     sql.NullInt32{Int32: int32(userAudioCount), Valid: true},
				}); err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	); err != nil {
		log.Println("Error in WriteInitialAudioRelCount transaction", err)
		return shared.InternalServerErrorDef()
	}

	return nil
}
