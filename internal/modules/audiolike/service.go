package audiolike

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type LikeUnlikeAudioParams struct {
	UserID  uuid.UUID
	AudioID uuid.UUID
}

func LikeAudio(
	ctx context.Context,
	db *database.Queries,
	params LikeUnlikeAudioParams,
) (*database.AudioLike, error) {
	err := sharedmodule.ValidateAudioExistsByID(ctx, db, params.AudioID)
	if err != nil {
		return nil, err
	}

	newAudioLike, err := db.CreateAudioLike(ctx, database.CreateAudioLikeParams{
		UserID:    params.UserID,
		AudioID:   params.AudioID,
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Println("Error creating audio like: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &newAudioLike, nil
}

func UnlikeAudio(
	ctx context.Context,
	db *database.Queries,
	params LikeUnlikeAudioParams,
) error {
	err := db.DeleteAudioLike(ctx, database.DeleteAudioLikeParams{
		UserID:  params.UserID,
		AudioID: params.AudioID,
	})
	if shared.IsDBErrorNotFound(err) {
		return shared.NotFound(shared.ErrAudioLikeNotFound)
	}

	if err != nil {
		log.Println("Error deleting audio like: ", err)
		return shared.InternalServerErrorDef()
	}

	return nil
}

func GetAudioLikesByUserID(
	ctx context.Context,
	db *database.Queries,
	userID uuid.UUID,
) ([]database.AudioLike, error) {
	entities, err := db.GetAudioLikesByUserID(ctx, userID)

	if err != nil {
		log.Println("Error getting audio likes by user ID:", err)
	}

	return entities, err
}

func GetAudioLikesByUserIDAndAudioIDs(
	ctx context.Context,
	db *database.Queries,
	params database.GetAudioLikesByUserIDAndAudioIDsParams,
) ([]database.AudioLike, error) {
	entities, err := db.GetAudioLikesByUserIDAndAudioIDs(ctx, params)

	if err != nil {
		log.Println("Error getting audio likes by user ID and audio IDs:", err)
	}

	return entities, err
}

func DeleteUserAudioLikesByAudioIDs(
	ctx context.Context,
	db *database.Queries,
	params database.DeleteUserAudioLikesByAudioIDsParams,
) error {
	err := db.DeleteUserAudioLikesByAudioIDs(ctx, params)

	if err != nil {
		log.Println("Error deleting user audio likes by audio IDs:", err)
	}

	return err
}
