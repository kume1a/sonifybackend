package playlist

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreatePlaylist(ctx context.Context, db *database.Queries, name string, thumbnailPath sql.NullString) (*database.Playlist, error) {
	entity, err := db.CreatePlaylist(ctx, database.CreatePlaylistParams{
		ID:            uuid.New(),
		Name:          name,
		ThumbnailPath: thumbnailPath,
	})

	if err != nil {
		log.Println("Error creating playlist:", err)
	}

	return &entity, err
}

func CreatePlaylistAudio(ctx context.Context, db *database.Queries, playlistID uuid.UUID, audioID uuid.UUID) (*database.PlaylistAudio, error) {
	entity, err := db.CreatePlaylistAudio(ctx, database.CreatePlaylistAudioParams{
		PlaylistID: playlistID,
		AudioID:    audioID,
	})

	if err != nil {
		log.Println("Error creating playlist audio:", err)
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

func GetPlaylistAudios(ctx context.Context, db *database.Queries, playlistID uuid.UUID) ([]database.GetPlaylistAudiosRow, error) {
	audios, err := db.GetPlaylistAudios(ctx, database.GetPlaylistAudiosParams{
		PlaylistID: playlistID,
	})

	if err != nil {
		log.Println("Error getting playlist audios:", err)
	}

	return audios, err
}
