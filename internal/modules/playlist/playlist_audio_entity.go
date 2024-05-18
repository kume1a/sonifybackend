package playlist

import (
	"context"
	"log"
	"time"

	"github.com/kume1a/sonifybackend/internal/database"
)

func CreatePlaylistAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreatePlaylistAudioParams,
) (*database.PlaylistAudio, error) {
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreatePlaylistAudio(ctx, params)

	if err != nil {
		log.Println("Error creating playlist audio:", err)
	}

	return &entity, err
}

func GetPlaylistAudioJoins(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudioJoinsParams,
) ([]database.GetPlaylistAudioJoinsRow, error) {
	audios, err := db.GetPlaylistAudioJoins(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audios:", err)
	}

	return audios, err
}

func DeletePlaylistAudiosByIds(
	ctx context.Context,
	db *database.Queries,
	params database.DeletePlaylistAudiosByIDsParams,
) error {
	err := db.DeletePlaylistAudiosByIDs(ctx, params)

	if err != nil {
		log.Println("Error deleting playlist audios by ids:", err)
	}

	return err
}

func GetPlaylistAudioJoinsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudioJoinsBySpotifyIDsParams,
) ([]database.GetPlaylistAudioJoinsBySpotifyIDsRow, error) {
	entities, err := db.GetPlaylistAudioJoinsBySpotifyIDs(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audio joins by spotify ids:", err)
	}

	return entities, err
}

func getPlaylistAudios(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudiosParams,
) ([]database.GetPlaylistAudiosRow, error) {
	audios, err := db.GetPlaylistAudios(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audios:", err)
	}

	return audios, err
}
