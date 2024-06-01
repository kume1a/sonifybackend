package userplaylist

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateUserPlaylist(
	ctx context.Context,
	db *database.Queries,
	params database.CreateUserPlaylistParams,
) (*database.UserPlaylist, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreateUserPlaylist(ctx, params)

	if err != nil {
		log.Println("Error creating user playlist:", err)
	}

	return &entity, err
}

func GetUserPlaylistsFull(
	ctx context.Context,
	db *database.Queries,
	params database.GetFullUserPlaylistsParams,
) ([]database.GetFullUserPlaylistsRow, error) {
	playlists, err := db.GetFullUserPlaylists(ctx, params)

	if err != nil {
		log.Println("Error getting user playlists full:", err)
	}

	return playlists, err
}

func GetUserPlaylistsByUserID(
	ctx context.Context,
	db *database.Queries,
	params database.GetUserPlaylistsParams,
) ([]database.UserPlaylist, error) {
	playlists, err := db.GetUserPlaylists(ctx, params)

	if err != nil {
		log.Println("Error getting user playlists:", err)
	}

	return playlists, err
}

func GetPlaylistIDsByUserID(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	playlistIds, err := db.GetPlaylistIDsByUserID(ctx, userId)

	if err != nil {
		log.Println("Error getting user playlist ids:", err)
	}

	return playlistIds, err
}

func GetUserPlaylistIDsByUserID(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	playlistIds, err := db.GetUserPlaylistIDsByUserID(ctx, userId)

	if err != nil {
		log.Println("Error getting user playlist ids:", err)
	}

	return playlistIds, err
}
