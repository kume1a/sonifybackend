package playlist

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreatePlaylist(
	ctx context.Context,
	db *database.Queries,
	params database.CreatePlaylistParams,
) (*database.Playlist, error) {
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

func GetPlaylists(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistsParams,
) ([]database.Playlist, error) {
	playlists, err := db.GetPlaylists(ctx, params)

	if err != nil {
		log.Println("Error getting playlists:", err)
	}

	return playlists, err
}

func GetSpotifyUserSavedPlaylistIds(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	playlistIds, err := db.GetSpotifyUserSavedPlaylistIDs(ctx, userId)

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

func DeletePlaylistsByIDs(
	ctx context.Context,
	db *database.Queries,
	playlistIds uuid.UUIDs,
) error {
	err := db.DeletePlaylistsByIDs(ctx, playlistIds)

	if err != nil {
		log.Println("Error deleting playlists by ids:", err)
	}

	return err
}

func DeletePlaylistByID(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) error {
	err := db.DeletePlaylistByID(ctx, playlistID)

	if err != nil {
		log.Println("Error deleting playlist by id:", err)
	}

	return err
}

func GetPlaylistByID(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) (*database.Playlist, error) {
	playlist, err := db.GetPlaylistByID(ctx, playlistID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, shared.NotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		log.Println("Error getting playlist by id:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &playlist, nil
}

func GetPlaylistBySpotifyID(
	ctx context.Context,
	db *database.Queries,
	spotifyID string,
) (*database.Playlist, error) {
	playlist, err := db.GetPlaylistBySpotifyID(ctx, spotifyID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, shared.NotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		log.Println("Error getting playlist by spotify id:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &playlist, nil
}

func GetPlaylistIDBySpotifyID(
	ctx context.Context,
	db *database.Queries,
	spotifyID string,
) (uuid.UUID, error) {
	playlistID, err := db.GetPlaylistIDBySpotifyID(ctx, spotifyID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return uuid.Nil, shared.NotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		return uuid.Nil, shared.InternalServerErrorDef()
	}

	return playlistID, nil
}

func UpdatePlaylistByID(
	ctx context.Context,
	db *database.Queries,
	params database.UpdatePlaylistByIDParams,
) (*database.Playlist, error) {
	playlist, err := db.UpdatePlaylistByID(ctx, params)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, shared.NotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		log.Println("Error updating playlist by id:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &playlist, err

}
