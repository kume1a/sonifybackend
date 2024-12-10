package userplaylist

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
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

type CreatePlaylistAndUserPlaylistParams struct {
	Name   string
	UserID uuid.UUID
}

func CreatePlaylistAndUserPlaylist(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params CreatePlaylistAndUserPlaylistParams,
) (*sharedmodule.UserPlaylistWithRel, error) {
	return shared.RunDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) (*sharedmodule.UserPlaylistWithRel, error) {
			playlist, err := playlist.CreatePlaylist(
				ctx, tx,
				database.CreatePlaylistParams{
					Name:              params.Name,
					AudioCount:        0,
					TotalAudioCount:   0,
					AudioImportStatus: database.ProcessStatusCOMPLETED,
				},
			)
			if err != nil {
				return nil, err
			}

			userPlaylist, err := CreateUserPlaylist(
				ctx, tx,
				database.CreateUserPlaylistParams{
					PlaylistID:             playlist.ID,
					UserID:                 params.UserID,
					IsSpotifySavedPlaylist: false,
				},
			)
			if err != nil {
				return nil, err
			}

			return &sharedmodule.UserPlaylistWithRel{
				UserPlaylist: userPlaylist,
				Playlist:     playlist,
			}, nil
		},
	)
}

type UpdateUserPlaylistParams struct {
	UserID         uuid.UUID
	UserPlaylistID uuid.UUID
	Name           string
}

func UpdateUserPlaylist(
	ctx context.Context,
	db *database.Queries,
	params UpdateUserPlaylistParams,
) (*sharedmodule.UserPlaylistWithRel, error) {
	userPlaylist, err := db.GetUserPlaylistByID(ctx, params.UserPlaylistID)
	if err != nil {
		log.Println("Error getting user playlist by ID:", err)

		if shared.IsDBErrorNotFound(err) {
			return nil, shared.NotFound(shared.ErrUserPlaylistNotFound)
		}

		return nil, err
	}

	entity, err := db.UpdatePlaylistByID(
		ctx,
		database.UpdatePlaylistByIDParams{
			PlaylistID: userPlaylist.PlaylistID,
			Name:       sql.NullString{String: params.Name, Valid: true},
		},
	)
	if err != nil {
		log.Println("Error updating user playlist:", err)
		return nil, err
	}

	return &sharedmodule.UserPlaylistWithRel{
		UserPlaylist: &userPlaylist,
		Playlist:     &entity,
	}, nil
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

func GetUserPlaylistUserIDsByPlaylistID(
	ctx context.Context,
	db *database.Queries,
	playlistId uuid.UUID,
) (uuid.UUIDs, error) {
	userIds, err := db.GetUserPlaylistUserIDsByPlaylistID(ctx, playlistId)

	if err != nil {
		log.Println("Error getting user playlist user ids:", err)
	}

	return userIds, err
}
