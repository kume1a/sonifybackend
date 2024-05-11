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
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

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

func GetPlaylistAudioJoins(ctx context.Context, db *database.Queries, playlistID uuid.UUID) ([]database.GetPlaylistAudioJoinsRow, error) {
	audios, err := db.GetPlaylistAudioJoins(ctx, database.GetPlaylistAudioJoinsParams{
		PlaylistID: playlistID,
	})

	if err != nil {
		log.Println("Error getting playlist audios:", err)
	}

	return audios, err
}

func CreateUserPlaylist(ctx context.Context, db *database.Queries, params database.CreateUserPlaylistParams) (*database.UserPlaylist, error) {
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreateUserPlaylist(ctx, params)

	if err != nil {
		log.Println("Error creating user playlist:", err)
	}

	return &entity, err
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

func GetUserPlaylists(
	ctx context.Context,
	db *database.Queries,
	params database.GetUserPlaylistsParams,
) ([]database.Playlist, error) {
	playlists, err := db.GetUserPlaylists(ctx, params)

	if err != nil {
		log.Println("Error getting user playlists:", err)
	}

	return playlists, err
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

func getPlaylistById(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) (*database.Playlist, error) {
	playlist, err := db.GetPlaylistById(ctx, playlistID)

	if err != nil {
		log.Println("Error getting playlist by id:", err)
	}

	return &playlist, err
}

func GetSpotifyUserSavedPlaylistIds(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	playlistIds, err := db.GetSpotifyUserSavedPlaylistIds(ctx, userId)

	if err != nil {
		log.Println("Error getting spotify user saved playlist ids:", err)
	}

	return playlistIds, err
}

func DeleteSpotifyUserSavedPlaylistJoins(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) error {
	err := db.DeleteSpotifyUserSavedPlaylistJoins(ctx, userId)

	if err != nil {
		log.Println("Error deleting spotify user saved playlist ids:", err)
	}

	return err
}

func DeletePlaylistsByIds(
	ctx context.Context,
	db *database.Queries,
	playlistIds uuid.UUIDs,
) error {
	err := db.DeletePlaylistsByIds(ctx, playlistIds)

	if err != nil {
		log.Println("Error deleting playlists by ids:", err)
	}

	return err
}

func GetUserPlaylistIDs(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	playlistIds, err := db.GetUserPlaylistIDs(ctx, userId)

	if err != nil {
		log.Println("Error getting user playlist ids:", err)
	}

	return playlistIds, err
}
