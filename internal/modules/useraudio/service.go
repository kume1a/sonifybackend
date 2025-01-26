package useraudio

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
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

	userAudio, err := db.CreateUserAudio(ctx, params)

	if err != nil {
		log.Println("Error creating user audio:", err)
	}

	if err := sharedmodule.IncrementAudioUserAudioCountByID(
		ctx, db, params.AudioID,
	); err != nil {
		log.Println("Error incrementing user audio count:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &userAudio, err
}

func BulkCreateUserAudiosTx(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params []database.CreateUserAudioParams,
) ([]database.UserAudio, error) {
	return shared.RunDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) ([]database.UserAudio, error) {
			userAudios := make([]database.UserAudio, 0, len(params))

			for _, param := range params {
				userAudio, err := CreateUserAudio(ctx, tx, param)
				if err != nil {
					log.Println("Error creating user audio:", err)
					return nil, shared.InternalServerErrorDef()
				}

				userAudios = append(userAudios, *userAudio)
			}

			return userAudios, nil
		},
	)
}

func UserAudioExistsByYoutubeVideoID(
	ctx context.Context,
	db *database.Queries,
	params database.UserAudioExistsByYoutubeVideoIDParams,
) (bool, error) {
	exists, err := db.UserAudioExistsByYoutubeVideoID(ctx, params)

	if err != nil {
		log.Println("Error checking if user audio exists by youtube video id: ", err)
	}

	return exists, err
}

func CountUserAudioByLocalID(
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

func GetUserAudioIDs(
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

func GetUserAudiosByAudioIDs(
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

func CountUserAudiosByAudioID(
	ctx context.Context,
	db *database.Queries,
	audioID uuid.UUID,
) (int64, error) {
	count, err := db.CountUserAudiosByAudioID(ctx, audioID)

	if err != nil {
		log.Println("Error counting user audios by audio id: ", err)
		return 0, shared.InternalServerErrorDef()
	}

	return count, err
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

func DeleteUserAudioTx(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params database.DeleteUserAudioParams,
) error {
	return shared.RunNoResultDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) error {
			return DeleteUserAudio(ctx, tx, params)
		},
	)
}

func DeleteUserAudio(
	ctx context.Context,
	tx *database.Queries,
	params database.DeleteUserAudioParams,
) error {
	count, err := tx.CountUserAudio(
		ctx,
		database.CountUserAudioParams(params),
	)

	if err != nil {
		log.Println("Error counting user audio: ", err)
		return err
	}

	if count == 0 {
		return shared.NotFound(shared.ErrUserAudioNotFound)
	}

	if err := tx.DeleteUserAudio(ctx, params); err != nil {
		log.Println("Error deleting user audio: ", err)
		return shared.InternalServerErrorDef()
	}

	if err := sharedmodule.DecrementAudioUserAudioCountByID(
		ctx, tx, params.AudioID,
	); err != nil {
		log.Println("Error decrementing user audio count: ", err)
		return shared.InternalServerErrorDef()
	}

	return nil
}
