package playlist

import (
	"context"
	"database/sql"
	"log"

	"github.com/kume1a/sonifybackend/internal/database"
)

func CreatePlaylist(ctx context.Context, db *database.Queries, name string, thumbnailPath sql.NullString) (*database.Playlist, error) {
	entity, err := db.CreatePlaylist(ctx, database.CreatePlaylistParams{
		Name:          name,
		ThumbnailPath: thumbnailPath,
	})

	if err != nil {
		log.Println("Error creating playlist:", err)
	}

	return &entity, err
}

func GetPlaylists(ctx context.Context, db *database.Queries, limit int32) ([]database.Playlist, error) {
	playlists, err := db.GetPlaylists(ctx, limit)

	if err != nil {
		log.Println("Error getting playlists:", err)
	}

	return playlists, err
}
