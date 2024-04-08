package playlistusecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func DeleteSpotifyUserSavedPlaylists(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	userId uuid.UUID,
) error {
	playlistIds, err := playlist.GetSpotifyUserSavedPlaylistIds(ctx, apiCfg.DB, userId)
	if err != nil {
		return err
	}

	if _, err := shared.RunDbTransaction(ctx, apiCfg, func(tx *database.Queries) (any, error) {
		err = playlist.DeleteSpotifyUserSavedPlaylistJoins(ctx, tx, userId)
		if err != nil {
			return nil, err
		}

		err = playlist.DeletePlaylistsByIds(ctx, tx, playlistIds)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}); err != nil {
		return err
	}

	return nil
}
