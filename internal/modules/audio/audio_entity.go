package audio

import (
	"context"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func createAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreateAudioParams,
) (*database.Audio, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}

	entity, err := db.CreateAudio(ctx, params)

	return &entity, err
}

func getAudioSpotifyIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) ([]database.GetAudioSpotifyIDsBySpotifyIDsRow, error) {
	return db.GetAudioSpotifyIDsBySpotifyIDs(ctx, spotifyIds)
}

func getAudioIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) (uuid.UUIDs, error) {
	return db.GetAudioIDsBySpotifyIDs(ctx, spotifyIds)
}
