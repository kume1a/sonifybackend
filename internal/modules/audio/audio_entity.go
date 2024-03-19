package audio

import (
	"context"
	"database/sql"
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

	audio, err := db.CreateAudio(ctx, database.CreateAudioParams{
		ID:             uuid.New(),
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
		Title:          title,
		Duration:       duration,
		Path:           sql.NullString{String: path, Valid: true},
		Author:         author,
		UserID:         uuid.NullUUID{UUID: userId, Valid: true},
		SizeBytes:      sizeBytes,
		YoutubeVideoID: youtubeVideoId,
		ThumbnailPath:  thumbnailPath,
	})

	return &audio, err
}

func GetUserAudioByYoutubeVideoId(
	db *database.Queries,
	ctx context.Context,
	userId uuid.UUID,
	youtubeVideoId string,
) (*database.Audio, error) {
	audio, err := db.GetUserAudioByVideoId(ctx, database.GetUserAudioByVideoIdParams{
		UserID:         uuid.NullUUID{UUID: userId, Valid: true},
		YoutubeVideoID: sql.NullString{String: youtubeVideoId, Valid: true},
	})

	return &audio, err
}
