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
) (*database.Audio, error) {
	createdAt := time.Now().UTC()

	audio, err := db.CreateAudio(ctx, database.CreateAudioParams{
		ID:        uuid.New(),
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Title:     title,
		Duration:  duration,
		Path:      sql.NullString{String: path, Valid: true},
		Author:    author,
		UserID:    uuid.NullUUID{UUID: userId, Valid: true},
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &audio, nil
}
