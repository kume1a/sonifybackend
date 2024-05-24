package audiolike

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateAudioLike(
	ctx context.Context,
	db *database.Queries,
	params database.CreateAudioLikeParams,
) (*database.AudioLike, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

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

func GetAudioLikes(
	ctx context.Context,
	db *database.Queries,
	params database.GetAudioLikesParams,
) ([]database.AudioLike, error) {
	entities, err := db.GetAudioLikes(ctx, params)

	if err != nil {
		log.Println("Error getting audio likes by user ID:", err)
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
