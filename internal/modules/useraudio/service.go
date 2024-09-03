package useraudio

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreateUserAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreateUserAudioParams,
) (*database.UserAudio, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreateUserAudio(ctx, params)

	if err != nil {
		log.Println("Error creating user audio:", err)
	}

	return &entity, err
}

func BulkCreateUserAudios(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params []database.CreateUserAudioParams,
) ([]database.UserAudio, error) {
	return shared.RunDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) ([]database.UserAudio, error) {
			audios := make([]database.UserAudio, 0, len(params))

			for _, param := range params {
				audio, err := CreateUserAudio(ctx, tx, param)
				if err != nil {
					log.Println("Error creating user audio:", err)
					return nil, shared.InternalServerErrorDef()
				}

				audios = append(audios, *audio)
			}

			return audios, nil
		},
	)
}

func GetUserAudioByYoutubeVideoId(
	ctx context.Context,
	db *database.Queries,
	params database.GetUserAudioByVideoIDParams,
) (*database.GetUserAudioByVideoIDRow, error) {
	audio, err := db.GetUserAudioByVideoID(ctx, params)

	if err != nil {
		log.Println("Error getting user audio by youtube video id: ", err)
	}

	return &audio, err
}

func CountUserAudioByLocalId(
	ctx context.Context,
	db *database.Queries,
	params database.CountUserAudioByLocalIDParams,
) (int64, error) {
	count, err := db.CountUserAudioByLocalID(ctx, params)

	if err != nil {
		log.Println("Error getting audio by local id: ", err)
	}

	return count, err
}

func GetUserAudioIds(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	ids, err := db.GetUserAudioIDs(ctx, userId)

	if err != nil {
		log.Println("Error getting user audio ids: ", err)
	}

	return ids, err
}

func GetUserAudiosByAudioIds(
	ctx context.Context,
	db *database.Queries,
	params database.GetUserAudiosByAudioIdsParams,
) ([]database.GetUserAudiosByAudioIdsRow, error) {
	if len(params.AudioIds) == 0 {
		return []database.GetUserAudiosByAudioIdsRow{}, nil
	}

	audios, err := db.GetUserAudiosByAudioIds(ctx, params)

	if err != nil {
		log.Println("Error getting audios by ids: ", err)
	}

	return audios, err
}

func UserAudioExists(
	ctx context.Context,
	db *database.Queries,
	params database.CountUserAudioParams,
) (bool, error) {
	count, err := db.CountUserAudio(ctx, params)

	if err != nil {
		log.Println("Error checking if user audio exists: ", err)
		return false, err
	}

	return count > 0, nil
}

func DeleteUserAudio(
	ctx context.Context,
	db *database.Queries,
	params database.DeleteUserAudioParams,
) error {
	count, err := db.CountUserAudio(ctx, database.CountUserAudioParams{
		UserID:  params.UserID,
		AudioID: params.AudioID,
	})

	if err != nil {
		log.Println("Error counting user audio: ", err)
		return err
	}

	if count == 0 {
		return shared.NotFound(shared.ErrUserAudioNotFound)
	}

	return db.DeleteUserAudio(ctx, params)
}
