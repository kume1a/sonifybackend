package playlist

import (
	"context"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func DeleteSpotifyUserSavedPlaylists(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	userId uuid.UUID,
) error {
	playlistIds, err := GetSpotifyUserSavedPlaylistIds(ctx, resourceConfig.DB, userId)
	if err != nil {
		return err
	}

	return shared.RunNoResultDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) error {
			err = DeleteSpotifyUserSavedPlaylistJoins(ctx, tx, userId)
			if err != nil {
				return err
			}

			err = DeletePlaylistsByIDs(ctx, tx, playlistIds)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func DeletePlaylistAndPlaylistAudiosBySpotifyID(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	playlistSpotifyID string,
) error {
	playlistId, err := GetPlaylistIDBySpotifyID(ctx, resourceConfig.DB, playlistSpotifyID)

	if err != nil && err.Error() == shared.ErrPlaylistNotFound {
		return nil
	}

	if err != nil {
		return err
	}

	return shared.RunNoResultDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) error {
			if err = playlistaudio.DeletePlaylistAudiosByPlaylistID(
				ctx, tx, playlistId,
			); err != nil {
				return err
			}

			if err := DeletePlaylistByID(ctx, tx, playlistId); err != nil {
				return err
			}

			return nil
		},
	)
}
