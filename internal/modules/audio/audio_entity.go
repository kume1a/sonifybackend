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
	db *database.Queries,
	ctx context.Context,
	title sql.NullString,
	author sql.NullString,
	duration sql.NullInt32,
	path string,
	userId uuid.UUID,
	sizeBytes sql.NullInt64,
	youtubeVideoId sql.NullString,
	thumbnailPath sql.NullString,
) (*database.Audio, error) {
	createdAt := time.Now().UTC()

	entity, err := db.CreateAudio(ctx, database.CreateAudioParams{
		ID:             uuid.New(),
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
		Title:          title,
		Duration:       duration,
		Path:           sql.NullString{String: path, Valid: true},
		Author:         author,
		SizeBytes:      sizeBytes,
		YoutubeVideoID: youtubeVideoId,
		ThumbnailPath:  thumbnailPath,
	})

	if err != nil {
		log.Println("Error creating audio:", err)
	}

	return &entity, err
}

func CreateUserAudio(db *database.Queries, ctx context.Context, userId uuid.UUID, audioId uuid.UUID) (*database.UserAudio, error) {
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
	db *database.Queries,
	ctx context.Context,
	userId uuid.UUID,
	youtubeVideoId string,
) (*database.GetUserAudioByVideoIdRow, error) {
	audio, err := db.GetUserAudioByVideoId(ctx, database.GetUserAudioByVideoIdParams{
		UserID:         userId,
		YoutubeVideoID: sql.NullString{String: youtubeVideoId, Valid: true},
	})

	return &audio, err
}
