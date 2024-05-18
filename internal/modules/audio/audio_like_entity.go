package audio

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateAudioLike(
	ctx context.Context,
	db *database.Queries,
	params database.CreateAudioLikeParams,
) (*database.AudioLike, error) {
	entity, err := db.CreateAudioLike(ctx, params)

	if err != nil {
		log.Println("Error creating audio like:", err)
	}

	return &entity, err
}

func DeleteAudioLike(
	ctx context.Context,
	db *database.Queries,
	params database.DeleteAudioLikeParams,
) error {
	err := db.DeleteAudioLike(ctx, params)

	if err != nil {
		log.Println("Error deleting audio like:", err)
	}

	return err
}

func GetAudioLikesByUserID(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) ([]database.AudioLike, error) {
	entities, err := db.GetAudioLikesByUserID(ctx, userId)

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
	if len(params.AudioIds) == 0 {
		return []database.AudioLike{}, nil
	}

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
