package hiddenuseraudio

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type HideUnhideAudioParams struct {
	UserID  uuid.UUID
	AudioID uuid.UUID
}

func HideUserAudio(
	ctx context.Context,
	db *database.Queries,
	params HideUnhideAudioParams,
) (*database.HiddenUserAudio, error) {
	err := sharedmodule.ValidateAudioExistsByID(ctx, db, params.AudioID)
	if err != nil {
		return nil, err
	}

	newHiddenUserAudio, err := db.CreateHiddenUserAudio(ctx, database.CreateHiddenUserAudioParams{
		UserID:    params.UserID,
		AudioID:   params.AudioID,
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Println("Error creating hidden user audio: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &newHiddenUserAudio, nil
}

func UnhideAudio(
	ctx context.Context,
	db *database.Queries,
	params HideUnhideAudioParams,
) error {
	err := db.DeleteHiddenUserAudio(ctx, database.DeleteHiddenUserAudioParams{
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

func GetHiddenUserAudiosByUserID(
	ctx context.Context,
	db *database.Queries,
	userID uuid.UUID,
) ([]database.HiddenUserAudio, error) {
	entities, err := db.GetHiddenUserAudiosByUserID(ctx, userID)

	if err != nil {
		log.Println("Error getting hidden user audios by user ID:", err)
	}

	return entities, err
}

func GetHiddenUserAudiosByUserIDAndAudioIDs(
	ctx context.Context,
	db *database.Queries,
	params database.GetHiddenUserAudiosByUserIDAndAudioIDsParams,
) ([]database.HiddenUserAudio, error) {
	entities, err := db.GetHiddenUserAudiosByUserIDAndAudioIDs(ctx, params)

	if err != nil {
		log.Println("Error getting hidden user audios by user ID and audio IDs:", err)
	}

	return entities, err
}
