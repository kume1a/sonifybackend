package playlist

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreatePlaylist(ctx context.Context, db *database.Queries, params database.CreatePlaylistParams) (*database.Playlist, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreatePlaylist(ctx, params)

	if err != nil {
		log.Println("Error creating playlist:", err)
	}

	return &entity, err
}

func CreatePlaylistAudio(ctx context.Context, db *database.Queries, params database.CreatePlaylistAudioParams) (*database.PlaylistAudio, error) {
	entity, err := db.CreatePlaylistAudio(ctx, params)

	if err != nil {
		log.Println("Error creating playlist audio:", err)
	}

	return &entity, err
}

func GetPlaylists(ctx context.Context, db *database.Queries, lastCreatedAt time.Time, limit int32) ([]database.Playlist, error) {
	playlists, err := db.GetPlaylists(ctx, database.GetPlaylistsParams{
		CreatedAt: lastCreatedAt,
		Limit:     limit,
	})

	if err != nil {
		log.Println("Error getting playlists:", err)
	}

	return playlists, err
}

func GetPlaylistAudios(ctx context.Context, db *database.Queries, playlistID uuid.UUID) ([]database.GetPlaylistAudiosRow, error) {
	audios, err := db.GetPlaylistAudios(ctx, database.GetPlaylistAudiosParams{
		PlaylistID: playlistID,
	})

	if err != nil {
		log.Println("Error getting playlist audios:", err)
	}

	return audios, err
}

func CreateUserPlaylist(ctx context.Context, db *database.Queries, params database.CreateUserPlaylistParams) (*database.UserPlaylist, error) {
	entity, err := db.CreateUserPlaylist(ctx, params)

	if err != nil {
		log.Println("Error creating user playlist:", err)
	}

	return &entity, err
}

func GetUserPlaylistsBySpotifyIds(ctx context.Context, db *database.Queries, params database.GetUserPlaylistsBySpotifyIdsParams) ([]database.Playlist, error) {
	playlists, err := db.GetUserPlaylistsBySpotifyIds(ctx, params)

	if err != nil {
		log.Println("Error getting user playlists by spotify ids:", err)
	}

	return playlists, err
}

func DeletePlaylistAudiosByIds(ctx context.Context, db *database.Queries, params database.DeletePlaylistAudiosByIdsParams) error {
	err := db.DeletePlaylistAudiosByIds(ctx, params)

	if err != nil {
		log.Println("Error deleting playlist audios by ids:", err)
	}

	return err
}

func GetPlaylistAudioJoinsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudioJoinsBySpotifyIdsParams,
) ([]database.GetPlaylistAudioJoinsBySpotifyIdsRow, error) {
	entities, err := db.GetPlaylistAudioJoinsBySpotifyIds(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audio joins by spotify ids:", err)
	}

	return entities, err
}
