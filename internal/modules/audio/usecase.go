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
	return shared.RunDBTransaction(
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

func DoesAudioExistByLocalId(ctx context.Context, db *database.Queries, userID uuid.UUID, localID string) (bool, error) {
	count, err := CountUserAudioByLocalId(ctx, db, database.CountUserAudioByLocalIdParams{
		LocalID: sql.NullString{String: localID, Valid: true},
		UserID:  userID,
	})
	if err != nil {
		if shared.IsDBErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return count > 0, nil
}

type WriteUserImportedLocalMusicParams struct {
	ApiConfig          *shared.ApiConfig
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
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	res, err := shared.RunDBTransaction(
		params.Context,
		params.ApiConfig,
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

			userAudio, err := CreateUserAudio(
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
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return res, nil
}

type SyncAudioLikesParams struct {
	Context  context.Context
	ApiCfg   *shared.ApiConfig
	UserID   uuid.UUID
	AudioIDs []uuid.UUID
}

func SyncAudioLikes(params SyncAudioLikesParams) error {
	audioLikes, err := GetAudioLikesByUserID(params.Context, params.ApiCfg.DB, params.UserID)
	if err != nil {
		return err
	}

	audioIDsSet := make(map[uuid.UUID]struct{}, len(params.AudioIDs))
	for _, audioID := range params.AudioIDs {
		audioIDsSet[audioID] = struct{}{}
	}

	audioLikesAudioIDsSet := make(map[uuid.UUID]struct{}, len(audioLikes))
	for _, audioLike := range audioLikes {
		audioLikesAudioIDsSet[audioLike.AudioID] = struct{}{}
	}

	toDelete := make([]uuid.UUID, 0)
	toInsert := make([]uuid.UUID, 0)

	for audioID := range audioLikesAudioIDsSet {
		if _, ok := audioIDsSet[audioID]; !ok {
			toDelete = append(toDelete, audioID)
		}
	}

	for audioID := range audioIDsSet {
		if _, ok := audioLikesAudioIDsSet[audioID]; !ok {
			toInsert = append(toInsert, audioID)
		}
	}

	shared.RunDBTransaction(
		params.Context,
		params.ApiCfg,
		func(tx *database.Queries) (any, error) {
			if err := DeleteUserAudioLikesByAudioIDs(params.Context, tx, database.DeleteUserAudioLikesByAudioIDsParams{
				UserID:   params.UserID,
				AudioIds: toDelete,
			}); err != nil {
				return nil, err
			}

			for _, audioID := range toInsert {
				if _, err := CreateAudioLike(params.Context, tx, database.CreateAudioLikeParams{
					UserID:  params.UserID,
					AudioID: audioID,
				}); err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	)

	return nil
}
