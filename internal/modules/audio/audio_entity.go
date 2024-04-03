package audio

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreateAudioParams,
) (*database.Audio, error) {
	createdAt := time.Now().UTC()

	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = createdAt
	}

	entity, err := db.CreateAudio(ctx, params)

	if err != nil {
		log.Println("Error creating audio:", err)
	}

	return &entity, err
}

func CreateUserAudio(ctx context.Context, db *database.Queries, userId uuid.UUID, audioId uuid.UUID) (*database.UserAudio, error) {
	entity, err := db.CreateUserAudio(ctx, database.CreateUserAudioParams{
		UserID:  userId,
		AudioID: audioId,
	})

	if err != nil {
		log.Println("Error creating user audio:", err)
	}

	return &entity, err
}

func GetUserAudioByYoutubeVideoId(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
	youtubeVideoId string,
) (*database.GetUserAudioByVideoIdRow, error) {
	audio, err := db.GetUserAudioByVideoId(ctx, database.GetUserAudioByVideoIdParams{
		UserID:         userId,
		YoutubeVideoID: sql.NullString{String: youtubeVideoId, Valid: true},
	})

	return &audio, err
}
