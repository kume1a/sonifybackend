package userplaylist

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateUserPlaylist(ctx context.Context, db *database.Queries, params database.CreateUserPlaylistParams) (*database.UserPlaylist, error) {
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
